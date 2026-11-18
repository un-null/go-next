package http

import (
	"net/http"
	"strconv"

	"backend/internal/entity"
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
	g.GET("/products/:id", h.GetProductByID)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	ctx := c.Request().Context()

	q := c.QueryParam("q")
	page := positiveNumber(c.QueryParam("page"), 1)
	limit := positiveNumber(c.QueryParam("limit"), 10)

	category := positiveNumber(c.QueryParam("category"), 0)

	var products []entity.Product
	var err error
	if category > 0 {
		products, err = h.productUseCase.GetProductsByCategory(ctx, category, q, page, limit)
	} else {
		products, err = h.productUseCase.GetAllProducts(ctx, q, page, limit)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}

func positiveNumber(raw string, fallback int) int {
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func (h *ProductHandler) GetProductByID(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	product, err := h.productUseCase.GetProductByID(ctx, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if product == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
	}

	return c.JSON(http.StatusOK, product)
}
