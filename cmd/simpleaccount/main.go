package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yeyee2901/test/config"
	"github.com/yeyee2901/test/internal/api"
	"github.com/yeyee2901/test/internal/logging"
	"github.com/yeyee2901/test/internal/utils"
)

// @title                   API Gateway - Simple Account
// @version                 1.0
// @BasePath                /
// @description.markdown

func main() {
	// load config
	cfg := config.MustLoadConfig("setting/setting.yaml")

	// setup logger
	logger := logging.NewFileLogger(cfg.Server.Logfile, cfg.Server.Name, slog.LevelInfo)
	slog.SetDefault(logger)

	db, err := connectDB(cfg)
	if err != nil {
		slog.Error("Cannot connect to database", "error", err)
		os.Exit(1)
	}

	server := api.NewAPIServer(cfg, db)

	server.RegisterMiddlewares()
	server.RegisterEndpoints()

	errChan := server.Run()
	slog.Info("Server is running")
	err = <-errChan
	if err != nil {
		slog.Error("Server exited")
		os.Exit(1)
	}
}

func connectDB(cfg *config.Config) (*sqlx.DB, error) {
	dsn := utils.BuildDatasourceName(utils.DataSource{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Database: cfg.DB.DBName,
	})
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to %s : %w", dsn, err)
	}

	return db, nil
}
