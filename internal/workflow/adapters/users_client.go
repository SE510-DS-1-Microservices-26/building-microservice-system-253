package adapters

import (
	"context"
	"fmt"
	"net/http"
)

type UsersClient struct {
	baseURL string
	client  *http.Client
}

func NewUsersClient(baseURL string) *UsersClient {
	return &UsersClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *UsersClient) ValidateUser(ctx context.Context, userID uint) error {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("users service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user %d not found", userID)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("users service returned %d for user %d", resp.StatusCode, userID)
	}

	return nil
}
