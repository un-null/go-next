package http

import (
	"net/http"
	"strconv"

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
	g.GET("/products/:id", h.GetProductById)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	products := h.productUseCase.GetAllProducts()
	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductById(c echo.Context) error {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	product := h.productUseCase.GetProductById(id)

	if product.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
	}

	return c.JSON(http.StatusOK, product)
}
