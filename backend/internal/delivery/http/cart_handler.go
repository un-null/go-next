package http

import (
	"net/http"
	"strconv"

	"backend/internal/entity"
	"backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	usecase *usecase.CartUseCase
}

func NewCartHandler(uc *usecase.CartUseCase) *CartHandler {
	return &CartHandler{usecase: uc}
}

func (h *CartHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/cart", h.AddToCart)
	g.GET("/cart", h.GetCartItems)
	g.DELETE("/cart", h.RemoveFromCart)
	g.DELETE("/cart/clear", h.ClearCart)
}

func (h *CartHandler) AddToCart(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, entity.AddToCartResponse{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	var request struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, entity.AddToCartResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if request.Quantity <= 0 {
		return c.JSON(http.StatusBadRequest, entity.AddToCartResponse{
			Success: false,
			Message: "Quantity must be greater than 0",
		})
	}

	err = h.usecase.AddToCart(c.Request().Context(), userID, request.ProductID, request.Quantity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.AddToCartResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.AddToCartResponse{
		Success: true,
		Message: "Product added to cart successfully",
	})
}

func (h *CartHandler) GetCartItems(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, entity.GetCartResponse{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	cartItems, err := h.usecase.GetCartItems(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.GetCartResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.GetCartResponse{
		Items:      cartItems,
		TotalItems: len(cartItems),
		Success:    true,
	})
}

func (h *CartHandler) RemoveFromCart(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, entity.RemoveFromCartResponse{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	productIDStr := c.QueryParam("product_id")
	if productIDStr == "" {
		return c.JSON(http.StatusBadRequest, entity.RemoveFromCartResponse{
			Success: false,
			Message: "product_id parameter is required",
		})
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, entity.RemoveFromCartResponse{
			Success: false,
			Message: "Invalid product_id format",
		})
	}

	err = h.usecase.RemoveFromCart(c.Request().Context(), userID, productID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.RemoveFromCartResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.RemoveFromCartResponse{
		Success: true,
		Message: "Product removed from cart successfully",
	})
}

func (h *CartHandler) ClearCart(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, entity.ClearCartResponse{
			Success: false,
			Message: "Invalid user ID",
		})
	}

	cartItems, err := h.usecase.GetCartItems(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ClearCartResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	itemsCount := len(cartItems)

	err = h.usecase.ClearCart(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ClearCartResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.ClearCartResponse{
		Success:      true,
		ItemsCleared: itemsCount,
		Message:      "Cart cleared successfully",
	})
}

func getUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userIDInterface := c.Get("userID")
	if userIDInterface == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found")
	}

	switch v := userIDInterface.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	default:
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid user ID type")
	}
}
