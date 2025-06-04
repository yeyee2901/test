package api

type DepositRequest struct {
	UserID int     `json:"user_id,omitempty" binding:"required"`
	Amount float64 `json:"amount,omitempty"  binding:"required"`
}

type DepositResponse struct {
	Status        string  `json:"status,omitempty"`
	TransactionID string  `json:"transaction_id,omitempty"`
	NewBalance    float64 `json:"new_balance,omitempty"`
}

type WithdrawRequest struct {
	UserID int     `json:"user_id,omitempty" binding:"required"`
	Amount float64 `json:"amount,omitempty"  binding:"required"`
}

type WithdrawResponse struct {
	UserID        int     `json:"user_id,omitempty"`
	TransactionID string  `json:"transaction_id,omitempty"`
	NewBalance    float64 `json:"new_balance,omitempty"`
}
