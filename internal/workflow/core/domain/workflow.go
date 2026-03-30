package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type WorkflowState string

const (
	WorkflowStateStarted       WorkflowState = "started"
	WorkflowStateUserValidated WorkflowState = "user_validated"
	WorkflowStateOrderCreated  WorkflowState = "order_created"
	WorkflowStateCompleted     WorkflowState = "completed"
	WorkflowStateCompensating  WorkflowState = "compensating"
	WorkflowStateCompensated   WorkflowState = "compensated"
	WorkflowStateFailed        WorkflowState = "failed"
)

type WorkflowType string

const (
	WorkflowTypePlaceOrder WorkflowType = "place-order"
)

var validTransitions = map[WorkflowState][]WorkflowState{
	WorkflowStateStarted:       {WorkflowStateUserValidated, WorkflowStateCompensating, WorkflowStateFailed},
	WorkflowStateUserValidated: {WorkflowStateOrderCreated, WorkflowStateCompensating, WorkflowStateFailed},
	WorkflowStateOrderCreated:  {WorkflowStateCompleted, WorkflowStateCompensating, WorkflowStateFailed},
	WorkflowStateCompensating:  {WorkflowStateCompensated, WorkflowStateFailed},
	WorkflowStateCompleted:     {},
	WorkflowStateCompensated:   {},
	WorkflowStateFailed:        {},
}

type WorkflowInstance struct {
	WorkflowID uuid.UUID       `json:"workflow_id"`
	Type       WorkflowType    `json:"type"`
	State      WorkflowState   `json:"state"`
	Payload    json.RawMessage `json:"payload"`
	LastError  *string         `json:"last_error"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func NewWorkflowInstance(wfType WorkflowType, payload json.RawMessage) *WorkflowInstance {
	return &WorkflowInstance{
		WorkflowID: uuid.New(),
		Type:       wfType,
		State:      WorkflowStateStarted,
		Payload:    payload,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (w *WorkflowInstance) Transition(next WorkflowState) error {
	for _, allowed := range validTransitions[w.State] {
		if allowed == next {
			w.State = next
			w.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("invalid transition: " + string(w.State) + " → " + string(next))
}

func (w *WorkflowInstance) IsTerminal() bool {
	return len(validTransitions[w.State]) == 0
}

func (w *WorkflowInstance) Fields() []interface{} {
	return []interface{}{
		&w.WorkflowID,
		&w.Type,
		&w.State,
		&w.Payload,
		&w.LastError,
		&w.CreatedAt,
		&w.UpdatedAt,
	}
}

type WorkflowTable struct{}

func (WorkflowTable) TableName() string {
	return "workflow_instances"
}

func (WorkflowTable) Columns() []string {
	return []string{
		"workflow_id",
		"type",
		"state",
		"payload",
		"last_error",
		"created_at",
		"updated_at",
	}
}
