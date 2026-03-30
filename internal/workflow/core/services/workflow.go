package services

import (
	"cafeteria-delivery/internal/workflow/core/domain"
	"cafeteria-delivery/internal/workflow/core/ports"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type WorkflowService struct {
	repo        ports.WorkflowRepository
	usersClient ports.UsersClientPort
	coreClient  ports.CoreClientPort
}

var _ ports.WorkflowService = (*WorkflowService)(nil)

func NewWorkflowService(
	repo ports.WorkflowRepository,
	usersClient ports.UsersClientPort,
	coreClient ports.CoreClientPort,
) *WorkflowService {
	return &WorkflowService{
		repo:        repo,
		usersClient: usersClient,
		coreClient:  coreClient,
	}
}

// StartPlaceOrder runs the place-order saga:
//  1. Validate user  → state: user_validated
//  2. Create order   → state: order_created
//  3. Complete       → state: completed
//
// On any failure after step 2 the compensation cancels the created order.
func (s *WorkflowService) StartPlaceOrder(ctx context.Context, req ports.PlaceOrderRequest) (*domain.WorkflowInstance, error) {
	// marshal request as saga payload for observability
	payload, _ := json.Marshal(req)

	wf := domain.NewWorkflowInstance(domain.WorkflowTypePlaceOrder, payload)

	if err := s.repo.Store(ctx, wf); err != nil {
		return nil, fmt.Errorf("failed to persist workflow: %w", err)
	}

	// ── Step 1: validate user ──────────────────────────────────────────────
	if err := s.usersClient.ValidateUser(ctx, req.UserID); err != nil {
		return s.fail(ctx, wf, fmt.Sprintf("step1 user validation failed: %v", err))
	}

	if err := s.transition(ctx, wf, domain.WorkflowStateUserValidated); err != nil {
		return nil, err
	}

	// ── Step 2: create order ───────────────────────────────────────────────
	items := make([]ports.CreateOrderItem, len(req.Items))
	for i, it := range req.Items {
		items[i] = ports.CreateOrderItem{ItemID: it.ItemID, Quantity: it.Quantity}
	}

	orderID, err := s.coreClient.CreateOrder(ctx, req.UserID, items)
	if err != nil {
		return s.fail(ctx, wf, fmt.Sprintf("step2 create order failed: %v", err))
	}

	if err = s.transition(ctx, wf, domain.WorkflowStateOrderCreated); err != nil {
		return nil, err
	}

	// ── Step 3: complete ───────────────────────────────────────────────────
	// (in a real system this could publish a RabbitMQ event, call another service, etc.)
	if err = s.transition(ctx, wf, domain.WorkflowStateCompleted); err != nil {
		// compensate: cancel the order that was already created
		s.compensate(ctx, wf, orderID, "step3 completion failed")
		return wf, nil
	}

	log.Printf("workflow %s completed successfully (order %d)", wf.WorkflowID, orderID)
	return wf, nil
}

func (s *WorkflowService) Find(ctx context.Context, id uuid.UUID) (*domain.WorkflowInstance, error) {
	return s.repo.Find(ctx, id)
}

// ── helpers ────────────────────────────────────────────────────────────────

// transition applies a state change and persists it
func (s *WorkflowService) transition(ctx context.Context, wf *domain.WorkflowInstance, next domain.WorkflowState) error {
	if err := wf.Transition(next); err != nil {
		return fmt.Errorf("state machine error: %w", err)
	}
	if err := s.repo.Update(ctx, wf); err != nil {
		return fmt.Errorf("failed to persist state %s: %w", next, err)
	}
	return nil
}

// fail records a terminal failure without compensation
func (s *WorkflowService) fail(ctx context.Context, wf *domain.WorkflowInstance, reason string) (*domain.WorkflowInstance, error) {
	log.Printf("workflow %s failed: %s", wf.WorkflowID, reason)
	_ = wf.Transition(domain.WorkflowStateFailed)
	wf.LastError = &reason
	_ = s.repo.Update(ctx, wf)
	return wf, nil
}

// compensate cancels the created order and marks the workflow as compensated
func (s *WorkflowService) compensate(ctx context.Context, wf *domain.WorkflowInstance, orderID uint, reason string) {
	log.Printf("workflow %s compensating: %s", wf.WorkflowID, reason)

	_ = wf.Transition(domain.WorkflowStateCompensating)
	wf.LastError = &reason
	_ = s.repo.Update(ctx, wf)

	if cancelErr := s.coreClient.CancelOrder(ctx, orderID); cancelErr != nil {
		log.Printf("workflow %s compensation failed to cancel order %d: %v", wf.WorkflowID, orderID, cancelErr)
		errMsg := fmt.Sprintf("%s; compensation error: %v", reason, cancelErr)
		_ = wf.Transition(domain.WorkflowStateFailed)
		wf.LastError = &errMsg
	} else {
		_ = wf.Transition(domain.WorkflowStateCompensated)
	}

	_ = s.repo.Update(ctx, wf)
}
