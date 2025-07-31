package http

import (
	"net/http"

	"backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: uc}
}

func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/users", h.GetUsers)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users := h.userUseCase.ListUsers()
	return c.JSON(http.StatusOK, users)
}
