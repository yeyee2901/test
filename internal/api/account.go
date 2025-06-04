package api

import "github.com/gin-gonic/gin"

// DepositRequest gin handler
// @Summary Adds balance to the account
// @Tags API
// @Param request body DepositRequest true "JSON body"
// @Produce json
// @Consume json
// @Success 200 {object} DepositResponse "Successful response"
// @Router /api/transactions/credit [post]
func (s *APIServer) DepositRequest(c *gin.Context) {}

// DepositRequest gin handler
// @Summary Adds balance to the account
// @Tags API
// @Param request body WithdrawRequest true "JSON body"
// @Produce json
// @Consume json
// @Success 200 {object} WithdrawResponse "Successful response"
// @Router /api/transactions/debit [post]
func (s *APIServer) WithdrawRequest(c *gin.Context) {}
