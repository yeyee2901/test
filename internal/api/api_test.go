package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yeyee2901/test/internal/account"
	"github.com/yeyee2901/test/internal/utils"
)

func TestAddBalance(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	testUsername := "auto_user_deposit_" + time.Now().Format("20060102150405")
	testBalance := 0.0
	testURL := "/test/deposit"

	testNumGoroutine := 100
	testTrxCount := 1 // per goroutine
	testAmount := 1000.00

	// PRECONDITION: connect to database
	db, err := connectDatabase()
	if err != nil {
		t.Fatal("precondition:", err)
	}

	// PRECONDITION: create test user
	err = createTestUser(db, testUsername, testBalance)
	if err != nil {
		t.Fatal("precondition:", err)
	}

	// PRECONDITION: create a simple server
	srv := gin.New()
	srv.Use(AttachRequestID())
	apiSrv := APIServer{
		db: db,
	}

	srv.Handle(http.MethodPost, testURL, apiSrv.DepositRequest)

	// PRECONDITION: construct the request
	body := DepositRequest{
		Username: testUsername,
		Amount:   testAmount,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	// TEST: fire up the goroutines
	wg := sync.WaitGroup{}
	wg.Add(testNumGoroutine)
	for i := 0; i < testNumGoroutine; i++ {
		go func(i int) {
			if t.Failed() {
				return
			}

			defer wg.Done()
			for count := 0; count < testTrxCount; count++ {
				req, err := http.NewRequest(http.MethodPost, testURL, bytes.NewBuffer(bodyJSON))
				if err != nil {
					t.Logf("goroutine #%d-%d failed with error %v", i, count, err)
					t.Fail()
				}

				resp := httptest.NewRecorder()
				srv.ServeHTTP(resp, req)

				t.Logf("goroutine #%d-%d : response HTTP %d", i, count, resp.Result().StatusCode)
				// body, _ := io.ReadAll(resp.Result().Body)
				// t.Log(string(body))
			}
		}(i)
	}

	wg.Wait()

	t.Log("done")

	// check user balance
	finalBalance := testBalance + (testAmount * float64(testNumGoroutine) * float64(testTrxCount))
	ewallet := account.NewSimpleEWalletSystem(db)
	user, err := ewallet.GetUser(testUsername)
	if err != nil {
		t.Fatal(err)
	}

	if user.Balance != finalBalance {
		t.Fatal("balance not equal. Want ", finalBalance, "; got ", user.Balance)
	}
}

func connectDatabase() (*sqlx.DB, error) {
	dsn := utils.BuildDatasourceName(utils.DataSource{
		User:     "postgres",
		Password: "your_password",
		Host:     "127.0.0.1:5432",
		Database: "simple_account",
	})

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// OPTIMIZE: fine tuning disini
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)
	return db, nil
}

func createTestUser(db *sqlx.DB, username string, balance float64) error {
	ewallet := account.NewSimpleEWalletSystem(db)
	return ewallet.CreateNewAccount(balance, username)
}
