package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cafeteria-delivery/internal/workflow/core/ports"
)

type CoreClient struct {
	baseURL string
	client  *http.Client
}

func NewCoreClient(baseURL string) *CoreClient {
	return &CoreClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type createOrderRequest struct {
	UserID uint              `json:"user_id"`
	Items  []createOrderItem `json:"items"`
}

type createOrderItem struct {
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
}

type createOrderResponse struct {
	ID uint `json:"id"`
}

func (c *CoreClient) CreateOrder(ctx context.Context, userID uint, items []ports.CreateOrderItem) (uint, error) {
	body, err := json.Marshal(createOrderRequest{UserID: userID, Items: func() []createOrderItem {
		out := make([]createOrderItem, len(items))
		for i, it := range items {
			out[i] = createOrderItem{ItemID: it.ItemID, Quantity: it.Quantity}
		}
		return out
	}()})
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/core/orders", bytes.NewReader(body))
	if err != nil {
		return 0, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("core service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("core service returned %d when creating order", resp.StatusCode)
	}

	var result createOrderResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.ID, nil
}

func (c *CoreClient) CancelOrder(ctx context.Context, orderID uint) error {
	url := fmt.Sprintf("%s/core/orders/%d", c.baseURL, orderID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("core service unavailable during compensation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("core service returned %d when cancelling order %d", resp.StatusCode, orderID)
	}

	return nil
}
