package http

import (
	"net/http"

	"backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{productUseCase: uc}
}

func (h *ProductHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/products", h.GetProducts)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	products := h.productUseCase.ListProducts()
	return c.JSON(http.StatusOK, products)
}
