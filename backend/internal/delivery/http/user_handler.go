package http

import (
	"net/http"
	"time"

	"backend/internal/entity"
	"backend/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	jwtSecret   string
}

func NewUserHandler(uc *usecase.UserUseCase, jwtSecret string) *UserHandler {
	return &UserHandler{userUseCase: uc, jwtSecret: jwtSecret}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users := h.userUseCase.ListUsers()
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/users", h.GetUsers)
	g.POST("/signup", h.SignUp)
	g.POST("/login", h.Login)
}

type signupRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *UserHandler) SignUp(c echo.Context) error {
	req := new(signupRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.userUseCase.SignUp(entity.User{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User created"})
}

type loginRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *UserHandler) Login(c echo.Context) error {
	req := new(loginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userUseCase.Login(req.Name, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := h.generateJWT(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
