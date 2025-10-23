package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/handlers"
	"github.com/tanydotai/tanyai/backend/internal/knowledge"
)

const defaultPort = "8080"

// Server wraps the Gin engine and exposes helpers for running the HTTP API.
type Server struct {
	engine     *gin.Engine
	httpServer *http.Server
}

// New constructs an HTTP server with all routes and middleware registered.
func New(database *sqlx.DB) *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	base := knowledge.LoadStaticKnowledgeBase()

	chatHandler := handlers.NewChatHandler(base)
	healthHandler := handlers.NewHealthHandler(database)

	engine.GET("/healthz", healthHandler.HandleHealth)

	api := engine.Group("/api/v1")
	{
		api.POST("/chat", chatHandler.HandleChat)
		api.GET("/knowledge-base", chatHandler.HandleKnowledgeBase)
	}

	httpSrv := &http.Server{
		Addr:         ":" + resolvePort(),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		engine:     engine,
		httpServer: httpSrv,
	}
}

// Run starts the HTTP server and blocks until shutdown is requested via context.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}

// Engine exposes the underlying Gin engine for testing.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func resolvePort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return defaultPort
	}
	return port
}
