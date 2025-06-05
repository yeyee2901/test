package api

type APIBaseResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type GetBalanceResponse struct {
	Balance float64 `json:"balance"`
}

type DepositRequest struct {
	Username string  `json:"username" binding:"required"`
	Amount   float64 `json:"amount"  binding:"required"`
}

type DepositResponse struct {
	APIBaseResponse
	TransactionID int     `json:"transaction_id,omitempty"`
	NewBalance    float64 `json:"new_balance,omitempty"`
}

type WithdrawRequest struct {
	Username string  `json:"username" binding:"required"`
	Amount   float64 `json:"amount"  binding:"required"`
}

type WithdrawResponse struct {
	APIBaseResponse
	UserID        int     `json:"user_id,omitempty"`
	TransactionID int     `json:"transaction_id,omitempty"`
	NewBalance    float64 `json:"new_balance,omitempty"`
}
