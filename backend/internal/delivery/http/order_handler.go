package http

import (
	"errors"
	"net/http"
	"strconv"

	"backend/internal/entity"
	"backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(uc usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: uc,
	}
}

func (h *OrderHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/orders", h.CreateOrder)
	g.GET("/orders/:id", h.GetOrderByID)
	g.GET("/users/:userID/orders", h.GetOrdersByUserID)
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	req := new(entity.CreateOrderRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	order, err := h.orderUsecase.CreateOrder(ctx, userID, *req)
	if err != nil {
		return h.respondError(c, err)
	}

	return c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	order, err := h.orderUsecase.GetOrderById(ctx, id)
	if err != nil {
		return h.respondError(c, err)
	}
	if order.UserID != userID {
		// Do not reveal that someone else's order exists.
		return c.JSON(http.StatusNotFound, map[string]string{"error": "order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrdersByUserID(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	pathUserID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}
	if pathUserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	orders, err := h.orderUsecase.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return h.respondError(c, err)
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) respondError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, entity.ErrInsufficientStock):
		return c.JSON(http.StatusConflict, map[string]string{"error": "insufficient stock"})
	case errors.Is(err, entity.ErrInsufficientCoins):
		return c.JSON(http.StatusPaymentRequired, map[string]string{"error": "insufficient coins"})
	case errors.Is(err, entity.ErrOrderNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": "order not found"})
	case errors.Is(err, entity.ErrProductNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
	case errors.Is(err, entity.ErrUserNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}
