package router

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
	"github.com/condratf/go-musthave-diploma-tpl/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrdersRepo struct {
	mock.Mock
}

func (m *mockOrdersRepo) Withdraw(login string, orderID string, amount float64) error {
	args := m.Called(login, orderID, amount)
	return args.Error(0)
}

func (m *mockOrdersRepo) UpdateOrderStatus(ctx context.Context, orderID string, status string, amount float64) error {
	args := m.Called(ctx, orderID, status, amount)
	return args.Error(0)
}

func (m *mockOrdersRepo) GetWithdrawals(login string) ([]repository.Withdrawal, error) {
	args := m.Called(login)
	return args.Get(0).([]repository.Withdrawal), args.Error(1)
}

func (m *mockOrdersRepo) GetPendingOrders(ctx context.Context) ([]models.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *mockOrdersRepo) GetOrders(login string) ([]models.Order, error) {
	args := m.Called(login)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *mockOrdersRepo) UploadOrder(login, order string) error {
	args := m.Called(login, order)
	return args.Error(0)
}

func (m *mockOrdersRepo) GetBalance(login string) (repository.BalanceResponse, error) {
	args := m.Called(login)
	return args.Get(0).(repository.BalanceResponse), args.Error(1)
}

type mockAppRouter struct {
	*AppRouter
	checkSessionFunc func(r *http.Request) (string, bool)
}

func (m *mockAppRouter) checkSession(r *http.Request) (string, bool) {
	return m.checkSessionFunc(r)
}

func TestUploadOrderHandler_Success(t *testing.T) {
	ordersRepo := new(mockOrdersRepo)
	ordersRepo.On("UploadOrder", "testuser", "123456").Return(nil)

	appRouter := &AppRouter{
		ordersRepo: ordersRepo,
	}

	// override checkSession
	wrapper := &mockAppRouter{
		AppRouter:        appRouter,
		checkSessionFunc: func(r *http.Request) (string, bool) { return "testuser", true },
	}

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString("123456"))
	w := httptest.NewRecorder()

	wrapper.uploadOrderHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	ordersRepo.AssertExpectations(t)
}
