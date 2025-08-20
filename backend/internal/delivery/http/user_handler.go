package http

import (
	"net/http"
	"time"

	"backend/internal/entity"
	"backend/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	jwtSecret   string
}

func NewUserHandler(uc *usecase.UserUseCase, jwtSecret string) *UserHandler {
	return &UserHandler{userUseCase: uc, jwtSecret: jwtSecret}
}

func (h *UserHandler) GetUserById(c echo.Context) error {
	idParam := c.Param("id")

	userID, err := uuid.Parse(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	user, err := h.userUseCase.GetUserById(c.Request().Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

func (h *UserHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/users/:id", h.GetUserById)
	g.POST("/signup", h.SignUp)
	g.POST("/login", h.Login)

	// Profile management routes (require authentication)
	protected := g.Group("")
	// protected.Use(h.AuthMiddleware()) // Add your auth middleware here
	protected.PATCH("/users/:id/name", h.UpdateUserName)
	protected.PATCH("/users/:id/email", h.UpdateUserEmail)
	protected.PATCH("/users/:id/password", h.ChangePassword)
	protected.PATCH("/users/:id/coins", h.UpdateUserCoins)
	protected.DELETE("/users/:id", h.DeleteUser)

	// Utility routes
	g.POST("/check-email", h.CheckEmailExists)
}

type signupRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
	Coins    int    `json:"coins,omitempty"`
}

func (h *UserHandler) SignUp(c echo.Context) error {
	req := new(signupRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userUseCase.SignUp(c.Request().Context(), entity.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Coins:    req.Coins,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"user":    user.ToResponse(),
	})
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
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

	user, err := h.userUseCase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := h.generateJWT(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user.ToResponse(),
	})
}

type updateNameRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

func (h *UserHandler) UpdateUserName(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	req := new(updateNameRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userUseCase.UpdateUserName(c.Request().Context(), userID, req.Name)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Name updated successfully",
		"user":    user.ToResponse(),
	})
}

type updateEmailRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (h *UserHandler) UpdateUserEmail(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	req := new(updateEmailRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userUseCase.UpdateUserEmail(c.Request().Context(), userID, req.Email)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Email updated successfully",
		"user":    user.ToResponse(),
	})
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	req := new(changePasswordRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.userUseCase.ChangePassword(c.Request().Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}

type updateCoinsRequest struct {
	Amount int `json:"amount" validate:"required"`
}

func (h *UserHandler) UpdateUserCoins(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	req := new(updateCoinsRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userUseCase.UpdateUserCoins(c.Request().Context(), userID, req.Amount)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Coins updated successfully",
		"user":    user.ToResponse(),
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	err = h.userUseCase.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

type checkEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (h *UserHandler) CheckEmailExists(c echo.Context) error {
	req := new(checkEmailRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	exists, err := h.userUseCase.CheckEmailExists(c.Request().Context(), req.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to check email")
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"exists": exists,
	})
}

func (h *UserHandler) parseUserID(c echo.Context) (uuid.UUID, error) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}
	return userID, nil
}

func (h *UserHandler) handleUseCaseError(err error) error {
	switch err.Error() {
	case "user not found":
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	case "email already exists":
		return echo.NewHTTPError(http.StatusConflict, "Email already exists")
	case "insufficient coins":
		return echo.NewHTTPError(http.StatusBadRequest, "Insufficient coins")
	case "current password is incorrect":
		return echo.NewHTTPError(http.StatusBadRequest, "Current password is incorrect")
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
}

func (h *UserHandler) generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Extended to 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
