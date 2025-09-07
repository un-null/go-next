package http

import (
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

	req := new(entity.CreateOrderRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	order, err := h.orderUsecase.CreateOrder(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	order, err := h.orderUsecase.GetOrderById(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if order == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrdersByUserID(c echo.Context) error {
	ctx := c.Request().Context()

	userIDParam := c.Param("userID")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	orders, err := h.orderUsecase.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}
