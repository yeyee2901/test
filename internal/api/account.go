package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeyee2901/test/internal/account"
)

// GetBalance gin handler
// @Summary Get balance info
// @Tags API
// @Param username query string false "Username"
// @Produce json
// @Success 200 {object} GetBalanceResponse "Successful response"
// @Success 400 {object} APIBaseResponse "Bad Request"
// @Success 404 {object} APIBaseResponse "User Not Found"
// @Success 500 {object} APIBaseResponse "Internal Server Error"
// @Router /api/balance [get]
func (s *APIServer) GetBalance(c *gin.Context) {
	logger := slog.Default().
		With(slog.String("request_id", c.GetString("X-Request-Id"))).
		With(slog.String("operation", "GetBalance"))

		// Validate the request
	username := c.Query("username")
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, APIBaseResponse{
			Status:  "error",
			Message: "Bad Request",
		})
		return
	}

	logger = logger.With(slog.Any("request_data", map[string]any{
		"username": username,
	}))

	ewallet := account.NewSimpleEWalletSystem(s.db)
	user, err := ewallet.GetUser(username)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrNotFound):
			logger.Error("failed to retrieve user", "error", err)
			c.AbortWithStatusJSON(http.StatusNotFound, APIBaseResponse{
				Status:  "error",
				Message: "user " + username + " not found",
			})

		default:
			logger.Error("failed to retrieve user", "error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, APIBaseResponse{
				Status:  "error",
				Message: "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, GetBalanceResponse{
		Balance: user.Balance,
	})
}

// DepositRequest gin handler
// @Summary Adds balance to the account
// @Tags API
// @Param request body DepositRequest true "JSON body"
// @Produce json
// @Consume json
// @Success 200 {object} DepositResponse "Successful response"
// @Success 400 {object} APIBaseResponse "Bad Request"
// @Success 500 {object} APIBaseResponse "Internal Server Error"
// @Router /api/transactions/credit [post]
func (s *APIServer) DepositRequest(c *gin.Context) {
	logger := slog.Default().
		With(slog.String("request_id", c.GetString("X-Request-Id"))).
		With(slog.String("operation", "DepositRequest"))

	// Validate the request
	req := new(DepositRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		logger.Error("validation failed on request", "error", err)

		c.AbortWithStatusJSON(http.StatusBadRequest, APIBaseResponse{
			Status:  "error",
			Message: "Bad Request",
		})
		return
	}

	ewallet := account.NewSimpleEWalletSystem(s.db)
	user, err := ewallet.GetUser(req.Username)
	if err != nil {
		logger.Error("failed to retrieve user", "error", err)

		switch {
		case errors.Is(err, account.ErrNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, APIBaseResponse{
				Status:  "error",
				Message: "user " + req.Username + " not found",
			})

		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, APIBaseResponse{
				Status:  "error",
				Message: "internal server error",
			})
		}

		return
	}

	trxResult, err := ewallet.AddBalance(user, req.Amount)
	if err != nil {
		logger.Error("failed to add balance", "error", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, APIBaseResponse{
			Status:  "error",
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, DepositResponse{
		APIBaseResponse: APIBaseResponse{
			Status: "success",
		},
		TransactionID: trxResult.ID,
		NewBalance:    user.Balance + req.Amount,
	})
}

// WithdrawRequest gin handler
// @Summary Deducts balance from the account
// @Tags API
// @Param request body WithdrawRequest true "JSON body"
// @Produce json
// @Consume json
// @Success 200 {object} WithdrawResponse "Successful response"
// @Success 400 {object} APIBaseResponse "Bad Request"
// @Success 500 {object} APIBaseResponse "Internal Server Error"
// @Router /api/transactions/debit [post]
func (s *APIServer) WithdrawRequest(c *gin.Context) {
	logger := slog.Default().
		With(slog.String("request_id", c.GetString("X-Request-Id"))).
		With(slog.String("operation", "WithdrawRequest"))

	// Validate the request
	req := new(DepositRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		logger.Error("validation failed on request", "error", err)

		c.AbortWithStatusJSON(http.StatusBadRequest, APIBaseResponse{
			Status:  "error",
			Message: "Bad Request",
		})
		return
	}

	ewallet := account.NewSimpleEWalletSystem(s.db)
	user, err := ewallet.GetUser(req.Username)
	if err != nil {
		logger.Error("failed to retrieve user", "error", err)
		switch {
		case errors.Is(err, account.ErrNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, APIBaseResponse{
				Status:  "error",
				Message: "user " + req.Username + " not found",
			})

		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, APIBaseResponse{
				Status:  "error",
				Message: "internal server error",
			})
		}

		return
	}

	trxResult, err := ewallet.DeductBalance(user, req.Amount)
	if err != nil {
		logger.Error("failed to deduct balance", "error", err)

		switch {
		case errors.Is(err, account.ErrInsufficient):
			c.AbortWithStatusJSON(http.StatusBadRequest, APIBaseResponse{
				Status:  "error",
				Message: "Insufficient funds",
			})

		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, APIBaseResponse{
				Status:  "error",
				Message: "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, DepositResponse{
		APIBaseResponse: APIBaseResponse{
			Status: "success",
		},
		TransactionID: trxResult.ID,
		NewBalance:    user.Balance - req.Amount,
	})
}
