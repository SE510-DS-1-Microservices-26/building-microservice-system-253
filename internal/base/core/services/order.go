package services

import (
	"cafeteria-delivery/internal/base/adapters"
	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/core/ports"
	"cafeteria-delivery/internal/base/dto"
	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/pkg/events"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	repo        ports.OrderRepository
	itemRepo    ports.ItemRepository
	usersClient *adapters.UsersClient
	publisher   *adapters.RabbitPublisher
}

var _ ports.OrderService = (*OrderService)(nil)

func NewOrderService(
	repository ports.OrderRepository,
	itemRepository ports.ItemRepository,
	usersClient *adapters.UsersClient,
	pub *adapters.RabbitPublisher,
) *OrderService {
	return &OrderService{
		repo:        repository,
		itemRepo:    itemRepository,
		usersClient: usersClient,
		publisher:   pub,
	}
}

// Store - store new order with its items
func (s *OrderService) Store(ctx context.Context, order *domain.Order) error {
	// validate, if user exists in users database
	if order.UserID != nil {
		if err := s.usersClient.Exists(ctx, *order.UserID); err != nil {
			return fmt.Errorf("failed to validate user: %w", err)
		}
	}

	// set order status to the pending
	order.Status = domain.OrderStatusPending

	// extract items ids
	ids := make([]uint, len(order.Items))
	for i, item := range order.Items {
		ids[i] = item.ItemID
	}

	// find menu items by ids
	items, err := s.itemRepo.FindByIDs(ctx, ids)
	if err != nil {
		return err
	}

	// create a map for items slice for faster lookup
	itemMap := make(map[uint]domain.Item, len(items))
	for _, it := range items {
		itemMap[it.ID] = it
	}

	for i, item := range order.Items {
		// check if the item is present from the order
		it, ok := itemMap[item.ItemID]
		if !ok {
			return fmt.Errorf("item %d: %w", item.ItemID, customErrors.ErrNotFound)
		}
		// check if item inventory has enough quantity for the order
		if it.Quantity < item.Quantity {
			return fmt.Errorf("insufficient quantity for item %d: %w", item.ItemID, customErrors.ErrConflict)
		}
		order.Items[i].UnitPrice = it.Price
	}

	// store the order with its items
	if err = s.repo.Store(ctx, order); err != nil {
		return err
	}

	// calculate the total price of the order
	for _, item := range order.Items {
		order.TotalPrice += item.UnitPrice * float64(item.Quantity)
	}

	// publish integration event
	summary := fmt.Sprintf("order #%d placed", order.ID)
	if order.UserID != nil {
		summary = fmt.Sprintf("order #%d placed by user %d", order.ID, *order.UserID)
	}
	payload := events.CoreItemCreatedEvent{
		EventID:       uuid.New().String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationID: uuid.New().String(),
		CoreItemID:    order.ID,
		OwnerUserID:   order.UserID,
		Summary:       summary,
	}
	publishErr := s.publisher.PublishOrderCreated(ctx, payload)
	if publishErr != nil {
		log.Printf("failed to publish core-item.created event for order %d: %v", order.ID, publishErr)
	}

	return nil
}

// Find - find the order with its items
func (s *OrderService) Find(ctx context.Context, id uint) (*domain.Order, error) {
	order, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	// calculate the total price for each order
	for _, item := range order.Items {
		order.TotalPrice += item.UnitPrice * float64(item.Quantity)
	}

	return order, nil
}

// List - list orders with their items
func (s *OrderService) List(ctx context.Context, filter dto.OrderFilter) ([]domain.Order, error) {
	orders, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// calculate the total price for each order
	for i, order := range orders {
		for _, item := range order.Items {
			orders[i].TotalPrice += item.UnitPrice * float64(item.Quantity)
		}
	}

	return orders, nil
}

// Count - count orders
func (s *OrderService) Count(ctx context.Context, filter dto.OrderFilter) (uint, error) {
	return s.repo.Count(ctx, filter)
}

// Update - update order
func (s *OrderService) Update(ctx context.Context, order *domain.Order) error {
	return s.repo.Update(ctx, order)
}

// Delete - delete order
func (s *OrderService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
