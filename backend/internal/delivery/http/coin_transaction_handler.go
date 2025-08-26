package http

import (
	"backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CoinTransactionHandler struct {
	coinTransactionUC usecase.CoinTransactionUseCase
}

func NewCoinTransactionHandler(coinTransactionUC usecase.CoinTransactionUseCase) *CoinTransactionHandler {
	return &CoinTransactionHandler{
		coinTransactionUC: coinTransactionUC,
	}
}

type chargeCoinsRequest struct {
	Amount      int    `json:"amount" validate:"required,gt=0"`
	Description string `json:"description" validate:"required"`
	OrderID     *int32 `json:"order_id,omitempty"`
}

type getTransactionsRequest struct {
	Page  int32 `query:"page" validate:"omitempty,gte=1"`
	Limit int32 `query:"limit" validate:"omitempty,gte=1,lte=100"`
}

type spendCoinsRequest struct {
	Amount      int    `json:"amount" validate:"required,gt=0"`
	Description string `json:"description" validate:"required"`
	OrderID     *int32 `json:"order_id,omitempty"`
}

func (h *CoinTransactionHandler) SpendUserCoins(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	req := new(spendCoinsRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, transaction, err := h.coinTransactionUC.SpendUserCoins(
		c.Request().Context(),
		userID,
		req.Amount,
		req.Description,
		req.OrderID,
	)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Coins spent successfully",
		"user":        user.ToResponse(),
		"transaction": transaction.ToResponse(),
	})
}

func (h *CoinTransactionHandler) ChargeUserCoins(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	req := new(chargeCoinsRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, transaction, err := h.coinTransactionUC.ChargeUserCoins(
		c.Request().Context(),
		userID,
		req.Amount,
		req.Description,
		req.OrderID,
	)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Coins charged successfully",
		"user":        user.ToResponse(),
		"transaction": transaction.ToResponse(),
	})
}

func (h *CoinTransactionHandler) GetUserTransactions(c echo.Context) error {
	userID, err := h.parseUserID(c)
	if err != nil {
		return err
	}

	req := new(getTransactionsRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	transactions, err := h.coinTransactionUC.GetUserTransactions(
		c.Request().Context(),
		userID,
		req.Page,
		req.Limit,
	)
	if err != nil {
		return h.handleUseCaseError(err)
	}

	response := make([]map[string]interface{}, len(transactions))
	for i, tx := range transactions {
		response[i] = tx.ToResponse()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": response,
		"page":         req.Page,
		"limit":        req.Limit,
	})
}

func (h *CoinTransactionHandler) GetTransactionByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid transaction ID")
	}

	transaction, err := h.coinTransactionUC.GetTransactionByID(c.Request().Context(), int32(id))
	if err != nil {
		return h.handleUseCaseError(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transaction": transaction.ToResponse(),
	})
}

func (h *CoinTransactionHandler) parseUserID(c echo.Context) (uuid.UUID, error) {
	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "User ID is not a string")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	return userID, nil
}

func (h *CoinTransactionHandler) handleUseCaseError(err error) error {
	switch err.Error() {
	case "user not found":
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case "insufficient coins":
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case "transaction not found":
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case "amount must be positive":
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case "description is required":
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
}
