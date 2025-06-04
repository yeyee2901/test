package account

import "github.com/jmoiron/sqlx"

type Account struct {
	UserID  int
	Balance float64
}

// AccountService is the interface responsible for operations
// on user Account
type AccountService interface {
	// CreateNewAccount creates a new account for the user with a initial balance
	CreateNewAccount(initBalance float64) (*Account, error)
}

type simpleAccountService struct {
	db *sqlx.DB
}
