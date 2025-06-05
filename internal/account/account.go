package account

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var ErrInsufficient = fmt.Errorf("account: insufficient funds")

const (
	TrxTypeCredit = "credit"
	TrxTypeDebit  = "debit"
)

type Account struct {
	ID        int          `db:"id"`
	Username  string       `db:"username"`
	Balance   float64      `db:"balance"`
	CreatedAt sql.NullTime `db:"created_at"`
}

type Transactions struct {
	ID        int          `db:"id"`
	UserID    int          `db:"user_id"`
	Amount    float64      `db:"amount"`
	TrxType   string       `db:"type"`
	CreatedAt sql.NullTime `db:"created_at"`
}

// EWalletSystem is the interface responsible for operations
// on user Account
type EWalletSystem interface {
	// CreateNewAccount creates a new account for the user with a initial balance
	CreateNewAccount(initBalance float64, userName string) error

	// GetUser gets user info with this username
	GetUser(username string) (*Account, error)

	// AddBalance adds fund for the user
	AddBalance(*Account, float64) error

	// DeductBalance deducts fund from the user
	DeductBalance(*Account, float64) error
}

type simpleEWallet struct {
	db *sqlx.DB
}

func NewSimpleEWalletSystem(db *sqlx.DB) EWalletSystem {
	return &simpleEWallet{
		db: db,
	}
}

// CreateNewAccount implements AccountService.
func (s *simpleEWallet) CreateNewAccount(initBalance float64, userName string) error {
	q := `
        INSERT INTO users
        (
            username,
            balance
        )
        VALUES 
        (
            :username,
            :balance
        )
    `

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(q, Account{
		Username: userName,
		Balance:  initBalance,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetUser implements AccountService.
func (s *simpleEWallet) GetUser(username string) (*Account, error) {
	q := `
        SELECT 
            id, username, balance, created_at
        FROM users
        WHERE username = $1
        LIMIT 1
    `

	acc := new(Account)
	err := s.db.Get(acc, q, username)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

// AddBalance implements AccountService.
func (s *simpleEWallet) AddBalance(acc *Account, amount float64) error {
	qBalance := `
        UPDATE users
        SET
            balance = balance + $1
        WHERE
            id = $2
    `
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(qBalance, amount, acc.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// if successful on adding balance, then create
	// the transaction record
	qTrx := `
        INSERT INTO transactions
            (user_id, amount, type)
        VALUES
            (:user_id, :amount, :type)
    `

	trx := Transactions{
		UserID:  acc.ID,
		Amount:  amount,
		TrxType: TrxTypeCredit,
	}

	_, err = tx.NamedExec(qTrx, trx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil
	}

	return nil
}

// DeductBalance implements EWalletSystem.
func (s *simpleEWallet) DeductBalance(acc *Account, amount float64) error {
	if !canDeductFund(acc, amount) {
		return ErrInsufficient
	}

	qBalance := `
        UPDATE users
        SET
            balance = balance - $1
        WHERE
            id = $2
    `
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(qBalance, amount, acc.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// if successful on adding balance, then create
	// the transaction record
	qTrx := `
        INSERT INTO transactions
            (user_id, amount, type)
        VALUES
            (:user_id, :amount, :type)
    `

	trx := Transactions{
		UserID:  acc.ID,
		Amount:  amount,
		TrxType: TrxTypeDebit,
	}

	_, err = tx.NamedExec(qTrx, trx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil
	}

	return nil
}

func canDeductFund(acc *Account, amount float64) bool {
	return (acc.Balance - amount) >= 0
}
