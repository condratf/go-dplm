package accrualservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

type AccrualService interface {
	GetOrderInfo(ctx context.Context, orderNumber string) (*models.AccrualResponse, error)
	RegisterOrder(ctx context.Context, orderNumber string) error
}

type accrualClient struct {
	baseURL string
	client  *http.Client
}

func NewAccrualClient(baseURL string, client *http.Client) AccrualService {
	return &accrualClient{
		baseURL: baseURL,
		client:  client,
	}
}

func (a *accrualClient) RegisterOrder(ctx context.Context, orderNumber string) error {
	url := fmt.Sprintf("%s/api/orders", a.baseURL)

	requestBody := struct {
		Order string `json:"order"`
	}{
		Order: orderNumber,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusAccepted:
		return nil
	case http.StatusConflict:
		return fmt.Errorf("order already registered")
	case http.StatusBadRequest:
		return fmt.Errorf("invalid order number format")
	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")
		return fmt.Errorf("too many requests, retry after %s", retryAfter)
	case http.StatusInternalServerError:
		return fmt.Errorf("accrual service internal error")
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (a *accrualClient) GetOrderInfo(ctx context.Context, orderNumber string) (*models.AccrualResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/orders/%s", a.baseURL, orderNumber), nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var result models.AccrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil
	case http.StatusNoContent:
		return nil, fmt.Errorf("order not registered")
	case http.StatusTooManyRequests:
		retryAfter := resp.Header.Get("Retry-After")
		return nil, fmt.Errorf("too many requests, retry after %s", retryAfter)
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("accrual service internal error")
	default:
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
}
