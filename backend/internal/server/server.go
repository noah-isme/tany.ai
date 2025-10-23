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
	adminhandlers "github.com/tanydotai/tanyai/backend/internal/handlers/admin"
	"github.com/tanydotai/tanyai/backend/internal/knowledge"
	"github.com/tanydotai/tanyai/backend/internal/middleware"
	"github.com/tanydotai/tanyai/backend/internal/repos"
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

	profileRepo := repos.NewProfileRepository(database)
	skillsRepo := repos.NewSkillRepository(database)
	servicesRepo := repos.NewServiceRepository(database)
	projectsRepo := repos.NewProjectRepository(database)

	profileHandler := adminhandlers.NewProfileHandler(profileRepo)
	skillHandler := adminhandlers.NewSkillHandler(skillsRepo)
	serviceHandler := adminhandlers.NewServiceHandler(servicesRepo)
	projectHandler := adminhandlers.NewProjectHandler(projectsRepo)
	uploadsHandler := adminhandlers.NewUploadsHandler()

	engine.GET("/healthz", healthHandler.HandleHealth)

	api := engine.Group("/api/v1")
	{
		api.POST("/chat", chatHandler.HandleChat)
		api.GET("/knowledge-base", chatHandler.HandleKnowledgeBase)
	}

	adminGroup := engine.Group("/api/admin", middleware.AuthzAdminStub())
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

		adminGroup.POST("/uploads", uploadsHandler.Create)
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
