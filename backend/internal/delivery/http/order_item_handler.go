package http

import (
	"net/http"
	"strconv"

	"backend/internal/entity"
	"backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type OrderItemHandler struct {
	orderItemUsecase usecase.OrderItemUsecase
}

func NewOrderItemHandler(uc usecase.OrderItemUsecase) *OrderItemHandler {
	return &OrderItemHandler{
		orderItemUsecase: uc,
	}
}

func (h *OrderItemHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/order-items", h.CreateOrderItem)
	g.GET("/order-items/:id", h.GetOrderItemByID)
	g.GET("/orders/:orderID/items", h.GetOrderItemsByOrderID)
}

func (h *OrderItemHandler) CreateOrderItem(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(entity.OrderItem)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	item, err := h.orderItemUsecase.CreateOrderItem(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *OrderItemHandler) GetOrderItemByID(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order item id"})
	}

	item, err := h.orderItemUsecase.GetOrderItemById(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if item == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "order item not found"})
	}

	return c.JSON(http.StatusOK, item)
}

func (h *OrderItemHandler) GetOrderItemsByOrderID(c echo.Context) error {
	ctx := c.Request().Context()

	orderIDParam := c.Param("orderID")
	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	items, err := h.orderItemUsecase.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, items)
}
