package http

import (
	"backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	usecase *usecase.CartUseCase
}

func NewCartHandler(uc *usecase.CartUseCase) *CartHandler {
	return &CartHandler{usecase: uc}
}

func (h *CartHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/cart", h.AddToCart)
	e.GET("/cart", h.GetCartItems)
	e.DELETE("/cart", h.RemoveFromCart)
	e.DELETE("/cart/clear", h.ClearCart)
}

func (h *CartHandler) AddToCart(c echo.Context) error {
	userID := c.Get("userID").(int)
	var request struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := h.usecase.AddToCart(userID, request.ProductID, request.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product added to cart"})
}

func (h *CartHandler) GetCartItems(c echo.Context) error {
	userID := c.Get("userID").(int)
	cartItems, err := h.usecase.GetCartItems(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, cartItems)
}

func (h *CartHandler) RemoveFromCart(c echo.Context) error {
	userID := c.Get("userID").(int)
	productIDStr := c.QueryParam("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product_id"})
	}

	err = h.usecase.RemoveFromCart(userID, productID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product removed from cart"})
}

func (h *CartHandler) ClearCart(c echo.Context) error {
	userID := c.Get("userID").(int)
	err := h.usecase.ClearCart(userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cart cleared"})
}
