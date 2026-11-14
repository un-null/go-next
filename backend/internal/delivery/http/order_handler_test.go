package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/entity"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrderUsecase struct {
	mock.Mock
}

func (m *mockOrderUsecase) GetOrderById(ctx context.Context, id int) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *mockOrderUsecase) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *mockOrderUsecase) CreateOrder(ctx context.Context, userID uuid.UUID, req entity.CreateOrderRequest) (*entity.Order, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

type stubValidator struct{}

func (stubValidator) Validate(i interface{}) error { return nil }

func newOrderContext(t *testing.T, method, path, body string, userID uuid.UUID) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	e.Validator = stubValidator{}

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID.String())
	return c, rec
}

func TestCreateOrder_TakesUserFromJWT(t *testing.T) {
	uc := new(mockOrderUsecase)
	userID := uuid.New()
	body := `{"items":[{"product_id":1,"quantity":2}],"user_id":"11111111-1111-1111-1111-111111111111","total_amount":1}`

	c, rec := newOrderContext(t, http.MethodPost, "/orders", body, userID)
	req := entity.CreateOrderRequest{Items: []entity.CreateOrderItemRequest{{ProductID: 1, Quantity: 2}}}
	uc.On("CreateOrder", mock.Anything, userID, req).Return(&entity.Order{ID: 1, UserID: userID}, nil)

	err := NewOrderHandler(uc).CreateOrder(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	uc.AssertExpectations(t)
}

func TestCreateOrder_InsufficientStock_Returns409(t *testing.T) {
	uc := new(mockOrderUsecase)
	userID := uuid.New()

	c, rec := newOrderContext(t, http.MethodPost, "/orders", `{"items":[{"product_id":1,"quantity":2}]}`, userID)
	uc.On("CreateOrder", mock.Anything, userID, mock.Anything).
		Return(nil, entity.ErrInsufficientStock)

	err := NewOrderHandler(uc).CreateOrder(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestGetOrderByID_OtherUser_Returns404(t *testing.T) {
	uc := new(mockOrderUsecase)
	userID := uuid.New()

	c, rec := newOrderContext(t, http.MethodGet, "/orders/1", "", userID)
	c.SetParamNames("id")
	c.SetParamValues("1")
	uc.On("GetOrderById", mock.Anything, 1).Return(&entity.Order{ID: 1, UserID: uuid.New()}, nil)

	err := NewOrderHandler(uc).GetOrderByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	uc.AssertExpectations(t)
}

func TestGetOrdersByUserID_OtherUser_Returns403(t *testing.T) {
	uc := new(mockOrderUsecase)
	userID := uuid.New()
	otherID := uuid.New()

	c, rec := newOrderContext(t, http.MethodGet, "/users/"+otherID.String()+"/orders", "", userID)
	c.SetParamNames("userID")
	c.SetParamValues(otherID.String())

	err := NewOrderHandler(uc).GetOrdersByUserID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	uc.AssertNotCalled(t, "GetOrdersByUserID", mock.Anything, mock.Anything)
}

func TestRespondError_MapsSentinels(t *testing.T) {
	uc := new(mockOrderUsecase)
	userID := uuid.New()
	h := NewOrderHandler(uc)

	cases := []struct {
		err    error
		status int
	}{
		{entity.ErrInsufficientStock, http.StatusConflict},
		{entity.ErrInsufficientCoins, http.StatusPaymentRequired},
		{entity.ErrOrderNotFound, http.StatusNotFound},
		{entity.ErrProductNotFound, http.StatusNotFound},
		{entity.ErrUserNotFound, http.StatusNotFound},
		{errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		c, rec := newOrderContext(t, http.MethodGet, "/orders/1", "", userID)
		err := h.respondError(c, tc.err)
		assert.NoError(t, err)
		assert.Equal(t, tc.status, rec.Code, tc.err.Error())
	}
}
