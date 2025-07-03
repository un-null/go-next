package http

import (
	"net/http"

	"backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	userUseCase *usecase.UserUseCase
}

func NewHandler(uc *usecase.UserUseCase) *Handler {
	return &Handler{userUseCase: uc}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/users", h.GetUsers)
}

func (h *Handler) GetUsers(c echo.Context) error {
	users := h.userUseCase.ListUsers()
	return c.JSON(http.StatusOK, users)
}
