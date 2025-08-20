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
	g.GET("/products/:id", h.GetProductByID)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	ctx := c.Request().Context()

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	products, err := h.productUseCase.GetAllProducts(ctx, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
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
