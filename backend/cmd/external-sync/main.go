package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/tanydotai/tanyai/backend/internal/config"
	"github.com/tanydotai/tanyai/backend/internal/db"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/tanydotai/tanyai/backend/internal/repos"
	"github.com/tanydotai/tanyai/backend/internal/services/ingest"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := db.Open(dbCtx, cfg.PostgresURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer database.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	sourceRepo := repos.NewExternalSourceRepository(database)
	itemRepo := repos.NewExternalItemRepository(database)
	ingestService := ingest.NewService(cfg.External.HTTPTimeout, cfg.External.RateLimitRPM, cfg.External.DomainAllowlist)

	if err := ensureDefaultSources(ctx, sourceRepo, cfg.External.SourcesDefault); err != nil {
		log.Fatalf("ensure defaults: %v", err)
	}

	if err := runSync(ctx, sourceRepo, itemRepo, ingestService); err != nil {
		log.Fatalf("sync failed: %v", err)
	}
}

func ensureDefaultSources(ctx context.Context, repo repos.ExternalSourceRepository, seeds []config.ExternalSourceSeed) error {
	defaults := make([]models.ExternalSource, 0, len(seeds))
	for _, seed := range seeds {
		name := strings.TrimSpace(seed.Name)
		baseURL := strings.TrimSpace(seed.BaseURL)
		if name == "" || baseURL == "" {
			continue
		}
		sourceType := strings.TrimSpace(seed.SourceType)
		if sourceType == "" {
			sourceType = "auto"
		}
		defaults = append(defaults, models.ExternalSource{
			Name:       name,
			BaseURL:    baseURL,
			SourceType: sourceType,
			Enabled:    seed.Enabled,
		})
	}
	if len(defaults) == 0 {
		return nil
	}
	return repo.EnsureDefaults(ctx, defaults)
}

func runSync(ctx context.Context, sourceRepo repos.ExternalSourceRepository, itemRepo repos.ExternalItemRepository, svc *ingest.Service) error {
	page := 1
	limit := 100
	synced := make([]syncResult, 0)

	for {
		params := repos.ListParams{Page: page, Limit: limit, SortField: "name", SortDir: "asc"}
		sources, total, err := sourceRepo.List(ctx, params)
		if err != nil {
			return err
		}
		for _, source := range sources {
			if !source.Enabled {
				continue
			}
			parsed, err := normalizeBaseURL(source.BaseURL)
			if err != nil {
				slog.Error("skip source", "id", source.ID, "error", err)
				continue
			}
			slog.Info("sync start", "id", source.ID, "name", source.Name)
			result, err := svc.Sync(ctx, ingest.Source{
				ID:           source.ID,
				Name:         source.Name,
				BaseURL:      parsed,
				SourceType:   source.SourceType,
				ETag:         source.ETag,
				LastModified: source.LastModified,
			})
			if err != nil {
				if errors.Is(err, ingest.ErrNotModified) {
					slog.Info("no changes", "id", source.ID)
					synced = append(synced, syncResult{ID: source.ID.String(), Name: source.Name, Status: "not_modified"})
					continue
				}
				slog.Error("sync failed", "id", source.ID, "error", err)
				synced = append(synced, syncResult{ID: source.ID.String(), Name: source.Name, Status: "error", Error: err.Error()})
				continue
			}

			if len(result.Items) > 0 {
				if err := itemRepo.Upsert(ctx, result.Items); err != nil {
					slog.Error("upsert items failed", "id", source.ID, "error", err)
					synced = append(synced, syncResult{ID: source.ID.String(), Name: source.Name, Status: "error", Error: err.Error()})
					continue
				}
			}
			if err := sourceRepo.UpdateSyncState(ctx, source.ID, result.ETag, result.LastModified, result.FetchedAt); err != nil {
				slog.Error("update sync state failed", "id", source.ID, "error", err)
				synced = append(synced, syncResult{ID: source.ID.String(), Name: source.Name, Status: "error", Error: err.Error()})
				continue
			}
			slog.Info("sync completed", "id", source.ID, "items", len(result.Items))
			synced = append(synced, syncResult{ID: source.ID.String(), Name: source.Name, Status: "ok", Items: len(result.Items)})
		}

		if int64(page*limit) >= total {
			break
		}
		page++
	}

	payload, err := json.MarshalIndent(struct {
		CompletedAt time.Time    `json:"completedAt"`
		Results     []syncResult `json:"results"`
	}{CompletedAt: time.Now(), Results: synced}, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(payload))
	return nil
}

type syncResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Items  int    `json:"items,omitempty"`
	Error  string `json:"error,omitempty"`
}

func normalizeBaseURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("base url required")
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "https://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	parsed.Fragment = ""
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed, nil
}
