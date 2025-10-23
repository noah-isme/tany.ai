package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tanydotai/tanyai/backend/internal/handlers"
	"github.com/tanydotai/tanyai/backend/internal/knowledge"
)

const defaultPort = "8080"

// Server wraps the Gin engine and exposes helpers for running the HTTP API.
type Server struct {
	engine *gin.Engine
	port   string
}

// New constructs an HTTP server with all routes and middleware registered.
func New() *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	base := knowledge.LoadStaticKnowledgeBase()

	chatHandler := handlers.NewChatHandler(base)
	healthHandler := handlers.NewHealthHandler()

	engine.GET("/healthz", healthHandler.HandleHealth)

	api := engine.Group("/api/v1")
	{
		api.POST("/chat", chatHandler.HandleChat)
		api.GET("/knowledge-base", chatHandler.HandleKnowledgeBase)
	}

	return &Server{
		engine: engine,
		port:   resolvePort(),
	}
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	return s.engine.Run(fmt.Sprintf(":%s", s.port))
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
