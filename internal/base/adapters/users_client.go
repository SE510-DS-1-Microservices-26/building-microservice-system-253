package adapters

import (
	"context"
	"fmt"
	"net/http"

	customErrors "cafeteria-delivery/internal/base/errors"
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

// Exists - checks if such user exists
func (c *UsersClient) Exists(ctx context.Context, userID uint) error {
	// construct a url
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	// create a new GET http request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to build http request: %w", err)
	}

	// send http request to users service
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call users service: %w", err)
	}
	defer resp.Body.Close()

	// response status code handling
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user %d: %w", userID, customErrors.ErrNotFound)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("users service returned %d for user %d", resp.StatusCode, userID)
	}

	return nil
}
