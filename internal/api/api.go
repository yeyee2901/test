package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/yeyee2901/test/config"
	"github.com/yeyee2901/test/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type APIConfig struct {
	Listener             string
	ServerTimeoutSeconds int
}

type APIServer struct {
	config *APIConfig

	gin        *gin.Engine
	db         *sqlx.DB
	httpServer *http.Server
}

func NewAPIServer(cfg *config.Config, db *sqlx.DB) *APIServer {
	if strings.ToLower(cfg.Server.Mode) == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &APIServer{
		config: &APIConfig{
			Listener:             cfg.Server.Listener,
			ServerTimeoutSeconds: cfg.Server.ServerTimeoutSeconds,
		},
		gin:        gin.New(),
		db:         db,
		httpServer: nil,
	}
}

func (api *APIServer) RegisterMiddlewares() {
	api.gin.Use(gin.Recovery())
	api.gin.Use(CORSMiddleware())
}

func (api *APIServer) RegisterEndpoints() {
	// Routes here...

	// register swagger
	docs.SwaggerInfo.Host = api.config.Listener
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	api.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Run runs the server. This will return an error channel that can
// be waited. This error channel will return a non-nil error
// whenever the server is stopped.
func (api *APIServer) Run() <-chan error {
	errChan := make(chan error)
	httpServer := &http.Server{
		Addr:         api.config.Listener,
		Handler:      api.gin,
		ReadTimeout:  time.Duration(api.config.ServerTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(api.config.ServerTimeoutSeconds) * time.Second,
	}

	go func() {
		fmt.Println("Server listening at:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	api.httpServer = httpServer

	return errChan
}

// Shutdown kills the HTTP Server entirely
func (s *APIServer) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}
