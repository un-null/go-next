package http

import (
	"net/http"
	"strconv"

	"backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
}

func NewCategoryHandler(uc *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{categoryUseCase: uc}
}

func (h *CategoryHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/categories", h.GetAllCategories)
	g.GET("/categories/:id", h.GetCategoryByID)
}

func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	ctx := c.Request().Context()

	categories, err := h.categoryUseCase.GetAllCategories(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category id"})
	}

	category, err := h.categoryUseCase.GetCategoryByID(ctx, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if category == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
	}

	return c.JSON(http.StatusOK, category)
}
