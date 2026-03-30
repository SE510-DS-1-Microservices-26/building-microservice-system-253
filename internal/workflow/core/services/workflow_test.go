package services

import (
	"context"
	"testing"

	"cafeteria-delivery/internal/workflow/core/domain"
	"cafeteria-delivery/internal/workflow/core/ports"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockWorkflowRepository struct {
	data map[uuid.UUID]*domain.WorkflowInstance
}

func newMockWorkflowRepository() *mockWorkflowRepository {
	return &mockWorkflowRepository{data: make(map[uuid.UUID]*domain.WorkflowInstance)}
}

func (m *mockWorkflowRepository) Store(_ context.Context, wf *domain.WorkflowInstance) error {
	m.data[wf.WorkflowID] = wf
	return nil
}

func (m *mockWorkflowRepository) Update(_ context.Context, wf *domain.WorkflowInstance) error {
	m.data[wf.WorkflowID] = wf
	return nil
}

func (m *mockWorkflowRepository) Find(_ context.Context, id uuid.UUID) (*domain.WorkflowInstance, error) {
	wf, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	return wf, nil
}

type mockUsersClient struct{ failUserID uint }

func (m *mockUsersClient) ValidateUser(_ context.Context, userID uint) error {
	if userID == m.failUserID {
		return assert.AnError
	}
	return nil
}

type mockCoreClient struct {
	failCreate   bool
	cancelCalled bool
}

func (m *mockCoreClient) CreateOrder(_ context.Context, _ uint, _ []ports.CreateOrderItem) (uint, error) {
	if m.failCreate {
		return 0, assert.AnError
	}
	return 42, nil
}

func (m *mockCoreClient) CancelOrder(_ context.Context, _ uint) error {
	m.cancelCalled = true
	return nil
}

func newTestService(u ports.UsersClientPort, c ports.CoreClientPort) *WorkflowService {
	return NewWorkflowService(newMockWorkflowRepository(), u, c)
}

func defaultRequest() ports.PlaceOrderRequest {
	return ports.PlaceOrderRequest{
		UserID: 1,
		Items:  []ports.PlaceOrderItem{{ItemID: 1, Quantity: 2}},
	}
}

func TestStartPlaceOrder_Success(t *testing.T) {
	svc := newTestService(&mockUsersClient{}, &mockCoreClient{})

	wf, err := svc.StartPlaceOrder(context.Background(), defaultRequest())

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Equal(t, domain.WorkflowStateCompleted, wf.State)
	assert.Equal(t, domain.WorkflowTypePlaceOrder, wf.Type)
	assert.Nil(t, wf.LastError)
	assert.True(t, wf.IsTerminal())
}

func TestStartPlaceOrder_UserValidationFails(t *testing.T) {
	svc := newTestService(&mockUsersClient{failUserID: 1}, &mockCoreClient{})

	wf, err := svc.StartPlaceOrder(context.Background(), defaultRequest())

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Equal(t, domain.WorkflowStateFailed, wf.State)
	assert.NotNil(t, wf.LastError)
	assert.Contains(t, *wf.LastError, "step1")
}

func TestStartPlaceOrder_CreateOrderFails_NoCompensation(t *testing.T) {
	coreClient := &mockCoreClient{failCreate: true}
	svc := newTestService(&mockUsersClient{}, coreClient)

	wf, err := svc.StartPlaceOrder(context.Background(), defaultRequest())

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Equal(t, domain.WorkflowStateFailed, wf.State)
	assert.Contains(t, *wf.LastError, "step2")
	assert.False(t, coreClient.cancelCalled)
}

func TestFind_ReturnsWorkflow(t *testing.T) {
	repo := newMockWorkflowRepository()
	svc := &WorkflowService{repo: repo, usersClient: &mockUsersClient{}, coreClient: &mockCoreClient{}}

	wf := domain.NewWorkflowInstance(domain.WorkflowTypePlaceOrder, []byte(`{}`))
	_ = repo.Store(context.Background(), wf)

	found, err := svc.Find(context.Background(), wf.WorkflowID)

	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, wf.WorkflowID, found.WorkflowID)
}

func TestStateMachine_InvalidTransition(t *testing.T) {
	wf := domain.NewWorkflowInstance(domain.WorkflowTypePlaceOrder, []byte(`{}`))

	err := wf.Transition(domain.WorkflowStateCompleted)

	assert.Error(t, err)
	assert.Equal(t, domain.WorkflowStateStarted, wf.State)
}

func TestStateMachine_HappyPath(t *testing.T) {
	wf := domain.NewWorkflowInstance(domain.WorkflowTypePlaceOrder, []byte(`{}`))

	require.NoError(t, wf.Transition(domain.WorkflowStateUserValidated))
	require.NoError(t, wf.Transition(domain.WorkflowStateOrderCreated))
	require.NoError(t, wf.Transition(domain.WorkflowStateCompleted))

	assert.True(t, wf.IsTerminal())
}

func TestStateMachine_CompensationPath(t *testing.T) {
	wf := domain.NewWorkflowInstance(domain.WorkflowTypePlaceOrder, []byte(`{}`))

	require.NoError(t, wf.Transition(domain.WorkflowStateUserValidated))
	require.NoError(t, wf.Transition(domain.WorkflowStateCompensating))
	require.NoError(t, wf.Transition(domain.WorkflowStateCompensated))

	assert.True(t, wf.IsTerminal())
}
