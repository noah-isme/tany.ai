package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/auth"
	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/handlers"
	adminhandlers "github.com/tanydotai/tanyai/backend/internal/handlers/admin"
	authhandlers "github.com/tanydotai/tanyai/backend/internal/handlers/auth"
	"github.com/tanydotai/tanyai/backend/internal/middleware"
	"github.com/tanydotai/tanyai/backend/internal/repos"
	"github.com/tanydotai/tanyai/backend/internal/services/kb"
	"github.com/tanydotai/tanyai/backend/internal/storage"
)

const defaultPort = "8080"

// Server wraps the Gin engine and exposes helpers for running the HTTP API.
type Server struct {
	engine     *gin.Engine
	httpServer *http.Server
}

// New constructs an HTTP server with all routes and middleware registered.
func New(database *sqlx.DB, cfg config.Config) (*Server, error) {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	aggregator := kb.NewAggregator(database, cfg.KnowledgeCacheTTL)
	chatHistoryRepo := repos.NewChatHistoryRepository(database)
	chatHandler := handlers.NewChatHandler(aggregator, chatHistoryRepo, cfg.ChatModel)
	healthHandler := handlers.NewHealthHandler(database)

	userRepo := repos.NewUserRepository(database)
	tokenService, err := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}
	rateLimiter := auth.NewRateLimiter(cfg.LoginRateLimitPerMin, cfg.LoginRateLimitBurst, 10*time.Minute)

	profileRepo := repos.NewProfileRepository(database)
	skillsRepo := repos.NewSkillRepository(database)
	servicesRepo := repos.NewServiceRepository(database)
	projectsRepo := repos.NewProjectRepository(database)

	profileHandler := adminhandlers.NewProfileHandler(profileRepo, aggregator.Invalidate)
	skillHandler := adminhandlers.NewSkillHandler(skillsRepo, aggregator.Invalidate)
	serviceHandler := adminhandlers.NewServiceHandler(servicesRepo, aggregator.Invalidate)
	projectHandler := adminhandlers.NewProjectHandler(projectsRepo, aggregator.Invalidate)

	objectStore, err := storage.New(cfg.Storage)
	if err != nil {
		return nil, err
	}
	uploadsLogger := log.New(os.Stdout, "", 0)
	uploadsHandler := adminhandlers.NewUploadsHandler(objectStore, cfg.Upload, uploadsLogger)
	uploadLimiter := auth.NewRateLimiter(cfg.UploadRateLimitPerMin, cfg.UploadRateLimitBurst, 10*time.Minute)
	authHandler := authhandlers.NewHandler(userRepo, tokenService, rateLimiter, cfg.RefreshCookieName)

	engine.GET("/healthz", healthHandler.HandleHealth)

	knowledgeLimiter := auth.NewRateLimiter(cfg.KnowledgeRateLimitPerMin, cfg.KnowledgeRateLimitBurst, 10*time.Minute)
	chatLimiter := auth.NewRateLimiter(cfg.ChatRateLimitPerMin, cfg.ChatRateLimitBurst, 10*time.Minute)

	api := engine.Group("/api/v1")
	{
		api.POST("/chat", middleware.RateLimitByIP(chatLimiter), middleware.JSONLogger("chat"), chatHandler.HandleChat)
		api.GET("/knowledge-base", middleware.RateLimitByIP(knowledgeLimiter), middleware.JSONLogger("knowledge_base"), chatHandler.HandleKnowledgeBase)
	}

	authGroup := engine.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
	}

	adminGroup := engine.Group("/api/admin", middleware.Authn(tokenService), middleware.AuthzAdmin())
	{
		adminGroup.GET("/profile", profileHandler.Get)
		adminGroup.PUT("/profile", profileHandler.Put)

		skills := adminGroup.Group("/skills")
		{
			skills.GET("", skillHandler.List)
			skills.POST("", skillHandler.Create)
			skills.PUT(":id", skillHandler.Update)
			skills.DELETE(":id", skillHandler.Delete)
			skills.PATCH("/reorder", skillHandler.Reorder)
		}

		services := adminGroup.Group("/services")
		{
			services.GET("", serviceHandler.List)
			services.POST("", serviceHandler.Create)
			services.PUT(":id", serviceHandler.Update)
			services.DELETE(":id", serviceHandler.Delete)
			services.PATCH("/reorder", serviceHandler.Reorder)
			services.PATCH(":id/toggle", serviceHandler.Toggle)
		}

		projects := adminGroup.Group("/projects")
		{
			projects.GET("", projectHandler.List)
			projects.POST("", projectHandler.Create)
			projects.PUT(":id", projectHandler.Update)
			projects.DELETE(":id", projectHandler.Delete)
			projects.PATCH("/reorder", projectHandler.Reorder)
			projects.PATCH(":id/feature", projectHandler.Feature)
		}

		adminGroup.POST("/uploads", middleware.RateLimitByIP(uploadLimiter), uploadsHandler.Create)
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
	}, nil
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
