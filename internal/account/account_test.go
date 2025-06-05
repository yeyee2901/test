package account

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yeyee2901/test/internal/utils"
)

func Test_simpleAccountService_AddBalance(t *testing.T) {
	testUsername := "test_user_" + time.Now().Format("20060102150405")
	testInitialBalance := 0.0
	db, err := connectDatabase()
	if err != nil {
		t.Fatal("precondition:", err)
	}
	ewallet := NewSimpleEWalletSystem(db)

	t.Cleanup(func() {
		err := deleteTestTrx(db)
		if err != nil {
			t.Log("post-test:", err)
		}

		err = deleteTestUsers(db)
		if err != nil {
			t.Log("post-test:", err)
		}
	})

	// create the account first
	err = ewallet.CreateNewAccount(testInitialBalance, testUsername)
	if err != nil {
		t.Fatal("precondition:", err)
	}

	acc, err := ewallet.GetUser(testUsername)
	if err != nil {
		t.Fatal("precondition:", err)
	}
	t.Log("Using account: ", acc.Username)

	// TEST: for errors
	tests := []struct {
		name    string
		acc     *Account
		amount  float64
		wantErr bool
	}{
		{
			name:    "test_success",
			acc:     acc,
			amount:  1200.0,
			wantErr: false,
		},
		{
			name:    "test_success",
			acc:     acc,
			amount:  2400.0,
			wantErr: false,
		},
		{
			name:    "test_success",
			acc:     acc,
			amount:  4800.0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ewallet.AddBalance(tt.acc, tt.amount)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddBalance() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddBalance() succeeded unexpectedly")
			}
		})
	}

	// TEST: final amount
	finalUser, err := ewallet.GetUser(testUsername)
	if err != nil {
		t.Fatal("failed to get user for final amount test", err)
	}

	finalAmount := 0.0
	for i := range tests {
		finalAmount += tests[i].amount
	}

	if finalUser.Balance != finalAmount {
		t.Fatal("Balance is incorrect. Want:", finalAmount, "; got:", finalUser.Balance)
	}
}

func Test_simpleAccountService_DeductBalance(t *testing.T) {
	testUsername := "test_user_" + time.Now().Format("20060102150405")
	testInitialBalance := 10000.0
	db, err := connectDatabase()
	if err != nil {
		t.Fatal("precondition:", err)
	}
	ewallet := NewSimpleEWalletSystem(db)

	t.Cleanup(func() {
		err := deleteTestTrx(db)
		if err != nil {
			t.Log("post-test:", err)
		}

		err = deleteTestUsers(db)
		if err != nil {
			t.Log("post-test:", err)
		}
	})

	// create the account first
	err = ewallet.CreateNewAccount(testInitialBalance, testUsername)
	if err != nil {
		t.Fatal("precondition:", err)
	}

	acc, err := ewallet.GetUser(testUsername)
	if err != nil {
		t.Fatal("precondition:", err)
	}
	t.Log("Using account: ", acc.Username)

	// TEST: for errors
	tests := []struct {
		name    string
		acc     *Account
		amount  float64
		wantErr bool
	}{
		{
			name:    "test_success",
			acc:     acc,
			amount:  1000.0,
			wantErr: false,
		},
		{
			name:    "test_success",
			acc:     acc,
			amount:  1000.0,
			wantErr: false,
		},
		{
			name:    "test_success",
			acc:     acc,
			amount:  1000.0,
			wantErr: false,
		},
		{
			name:    "test_insufficient",
			acc:     acc,
			amount:  100000.00,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := ewallet.DeductBalance(tt.acc, tt.amount)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddBalance() failed: %v", gotErr)
				}
				t.Log("log: got error", gotErr)
				return
			}
			if tt.wantErr {
				t.Fatal("AddBalance() succeeded unexpectedly")
			}
		})
	}

	// TEST: final amount
	finalUser, err := ewallet.GetUser(testUsername)
	if err != nil {
		t.Fatal("failed to get user for final amount test", err)
	}

	deductTotal := 0.0
	for i := range tests {
		// skip amounts that fail
		if tests[i].wantErr {
			continue
		}
		deductTotal += tests[i].amount
	}
	finalAmount := testInitialBalance - deductTotal

	if finalUser.Balance != finalAmount {
		t.Fatal("Balance is incorrect. Want:", finalAmount, "; got:", finalUser.Balance)
	}
}

func connectDatabase() (*sqlx.DB, error) {
	dsn := utils.BuildDatasourceName(utils.DataSource{
		User:     "postgres",
		Password: "your_password",
		Host:     "127.0.0.1:5432",
		Database: "simple_account",
	})

	return sqlx.Connect("postgres", dsn)
}

func deleteTestUsers(db *sqlx.DB) error {
	q := `DELETE FROM users WHERE username LIKE '%test_user%'`
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(q)
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

func deleteTestTrx(db *sqlx.DB) error {
	q := `
        DELETE FROM transactions trx 
        USING users u
        WHERE u.username LIKE '%test_user%' AND u.id = trx.user_id
    `
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(q)
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
