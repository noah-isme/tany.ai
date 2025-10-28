package ingest

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tanydotai/tanyai/backend/internal/models"
	"github.com/temoto/robotstxt"
	"golang.org/x/time/rate"
)

const (
	defaultUserAgent       = "tanyai-external-sync/1.0"
	defaultMaxPages        = 64
	maxBodyBytes     int64 = 8 << 20 // 8 MiB safety limit
	defaultRetries         = 3
)

type conditionalHeaders struct {
	etag         *string
	lastModified *time.Time
}

// Service orchestrates syncing external knowledge sources.
type Service struct {
	client     *http.Client
	limiter    *rate.Limiter
	sanitizer  *bluemonday.Policy
	allowlist  map[string]struct{}
	robotsMu   sync.Mutex
	robots     map[string]*robotstxt.RobotsData
	userAgent  string
	maxPages   int
	retryCount int
}

// Option configures the Service.
type Option func(*Service)

// WithUserAgent overrides the default user agent string.
func WithUserAgent(agent string) Option {
	return func(s *Service) {
		if strings.TrimSpace(agent) != "" {
			s.userAgent = agent
		}
	}
}

// WithMaxPages sets a cap on how many pages will be crawled per sync.
func WithMaxPages(max int) Option {
	return func(s *Service) {
		if max > 0 {
			s.maxPages = max
		}
	}
}

// WithRetry overrides the default retry attempts.
func WithRetry(retries int) Option {
	return func(s *Service) {
		if retries >= 0 {
			s.retryCount = retries
		}
	}
}

// NewService builds a Service with sane defaults.
func NewService(timeout time.Duration, rpm int, allowlist []string, opts ...Option) *Service {
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	if rpm <= 0 {
		rpm = 30
	}
	interval := time.Minute / time.Duration(rpm)
	if interval <= 0 {
		interval = time.Second
	}
	allowed := make(map[string]struct{}, len(allowlist))
	for _, host := range allowlist {
		host = strings.ToLower(strings.TrimSpace(host))
		if host != "" {
			allowed[host] = struct{}{}
		}
	}

	svc := &Service{
		client:     &http.Client{Timeout: timeout},
		limiter:    rate.NewLimiter(rate.Every(interval), 1),
		sanitizer:  bluemonday.StrictPolicy(),
		allowlist:  allowed,
		robots:     make(map[string]*robotstxt.RobotsData),
		userAgent:  defaultUserAgent,
		maxPages:   defaultMaxPages,
		retryCount: defaultRetries,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Sync retrieves normalized items for a source. It performs conditional requests
// using the stored ETag/Last-Modified metadata.
func (s *Service) Sync(ctx context.Context, source Source) (Result, error) {
	if source.BaseURL == nil {
		return Result{}, errors.New("base url missing")
	}
	if err := s.ensureHostAllowed(source.BaseURL.Host); err != nil {
		return Result{}, err
	}

	sitemap, err := s.fetchSitemap(ctx, source)
	if err != nil {
		return Result{}, err
	}

	items, err := s.fetchPages(ctx, source, sitemap.URLs)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Items:        items,
		ETag:         sitemap.ETag,
		LastModified: sitemap.LastModified,
		FetchedAt:    time.Now(),
	}, nil
}

type sitemapPayload struct {
	URLs         []string
	ETag         *string
	LastModified *time.Time
}

func (s *Service) fetchSitemap(ctx context.Context, source Source) (sitemapPayload, error) {
	candidates := []string{"/sitemap-index.xml", "/sitemap.xml"}
	cond := &conditionalHeaders{etag: source.ETag, lastModified: source.LastModified}

	for _, path := range candidates {
		target := source.BaseURL.ResolveReference(&url.URL{Path: path})
		payload, err := s.retrieveSitemap(ctx, target, cond)
		if err == nil {
			return payload, nil
		}
		if errors.Is(err, ErrNotModified) {
			return sitemapPayload{}, ErrNotModified
		}
		// Try next candidate
	}
	return sitemapPayload{}, fmt.Errorf("no sitemap available for %s", source.BaseURL)
}

func (s *Service) retrieveSitemap(ctx context.Context, target *url.URL, cond *conditionalHeaders) (sitemapPayload, error) {
	resp, body, meta, err := s.fetch(ctx, target, cond, false)
	if err != nil {
		return sitemapPayload{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return sitemapPayload{}, ErrNotModified
	}
	if resp.StatusCode >= 400 {
		return sitemapPayload{}, fmt.Errorf("fetch sitemap %s: status %d", target, resp.StatusCode)
	}

	urls, err := s.parseSitemap(ctx, target, body)
	if err != nil {
		return sitemapPayload{}, err
	}
	if len(urls) == 0 {
		return sitemapPayload{}, fmt.Errorf("sitemap empty for %s", target)
	}

	result := sitemapPayload{URLs: urls}
	if meta.etag != nil {
		result.ETag = meta.etag
	}
	if meta.lastModified != nil {
		result.LastModified = meta.lastModified
	}
	return result, nil
}

func (s *Service) parseSitemap(ctx context.Context, base *url.URL, body []byte) ([]string, error) {
	type urlEntry struct {
		Loc string `xml:"loc"`
	}
	type sitemapIndex struct {
		Sitemaps []urlEntry `xml:"sitemap"`
	}
	type urlSet struct {
		URLs []urlEntry `xml:"url"`
	}

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.Strict = false

	var root xml.Name
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if start, ok := token.(xml.StartElement); ok {
			root = start.Name
			break
		}
	}

	decoder = xml.NewDecoder(bytes.NewReader(body))
	decoder.Strict = false

	var urls []string
	switch root.Local {
	case "sitemapindex":
		var index sitemapIndex
		if err := decoder.Decode(&index); err != nil {
			return nil, err
		}
		for _, entry := range index.Sitemaps {
			loc := strings.TrimSpace(entry.Loc)
			if loc == "" {
				continue
			}
			u, err := url.Parse(loc)
			if err != nil {
				continue
			}
			if !u.IsAbs() {
				u = base.ResolveReference(u)
			}
			if err := s.ensureHostAllowed(u.Host); err != nil {
				continue
			}
			nested, err := s.retrieveNestedSitemap(ctx, u)
			if err != nil {
				continue
			}
			urls = append(urls, nested...)
			if len(urls) >= s.maxPages {
				return urls[:s.maxPages], nil
			}
		}
	case "urlset":
		var set urlSet
		if err := decoder.Decode(&set); err != nil {
			return nil, err
		}
		for _, entry := range set.URLs {
			loc := strings.TrimSpace(entry.Loc)
			if loc == "" {
				continue
			}
			u, err := url.Parse(loc)
			if err != nil {
				continue
			}
			if !u.IsAbs() {
				u = base.ResolveReference(u)
			}
			if err := s.ensureHostAllowed(u.Host); err != nil {
				continue
			}
			urls = append(urls, u.String())
			if len(urls) >= s.maxPages {
				return urls[:s.maxPages], nil
			}
		}
	default:
		return nil, fmt.Errorf("unsupported sitemap root %s", root.Local)
	}

	return urls, nil
}

func (s *Service) retrieveNestedSitemap(ctx context.Context, target *url.URL) ([]string, error) {
	resp, body, _, err := s.fetch(ctx, target, nil, false)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("nested sitemap %s returned %d", target, resp.StatusCode)
	}
	return s.parseSitemap(ctx, target, body)
}

func (s *Service) fetchPages(ctx context.Context, source Source, urls []string) ([]models.ExternalItem, error) {
	seen := make(map[string]struct{})
	items := make([]models.ExternalItem, 0, len(urls))

	for _, raw := range urls {
		if len(items) >= s.maxPages {
			break
		}
		pageURL, err := url.Parse(raw)
		if err != nil {
			continue
		}
		if err := s.ensureHostAllowed(pageURL.Host); err != nil {
			continue
		}
		resp, body, _, err := s.fetch(ctx, pageURL, nil, false)
		if err != nil {
			continue
		}
		func() {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return
			}
			normalized := s.extractItems(source, pageURL, body)
			for _, item := range normalized {
				if _, exists := seen[item.Hash]; exists {
					continue
				}
				seen[item.Hash] = struct{}{}
				items = append(items, item)
				if len(items) >= s.maxPages {
					break
				}
			}
		}()
	}

	return items, nil
}

func (s *Service) extractItems(source Source, pageURL *url.URL, body []byte) []models.ExternalItem {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil
	}
	var results []models.ExternalItem
	doc.Find("script[type='application/ld+json']").Each(func(_ int, sel *goquery.Selection) {
		payload := strings.TrimSpace(sel.Text())
		if payload == "" {
			return
		}
		dec := json.NewDecoder(strings.NewReader(payload))
		dec.UseNumber()

		var data any
		if err := dec.Decode(&data); err != nil {
			return
		}
		switch value := data.(type) {
		case []any:
			for _, entry := range value {
				if item := s.itemFromLD(source, pageURL, entry); item != nil {
					results = append(results, *item)
				}
			}
		default:
			if item := s.itemFromLD(source, pageURL, value); item != nil {
				results = append(results, *item)
			}
		}
	})
	return results
}

func (s *Service) itemFromLD(source Source, pageURL *url.URL, entry any) *models.ExternalItem {
	obj, ok := entry.(map[string]any)
	if !ok {
		return nil
	}

	typeVal := stringFrom(obj["@type"])
	title := strings.TrimSpace(firstNonEmpty(
		stringFrom(obj["name"]),
		stringFrom(obj["headline"]),
		stringFrom(obj["title"]),
	))
	if title == "" {
		return nil
	}

	summary := strings.TrimSpace(firstNonEmpty(
		stringFrom(obj["headline"]),
		stringFrom(obj["description"]),
	))
	content := strings.TrimSpace(firstNonEmpty(
		stringFrom(obj["about"]),
		stringFrom(obj["articleBody"]),
	))

	publishedAt := parseTime(stringFrom(obj["datePublished"]))

	metadata := models.JSONB{}
	if image := stringFrom(obj["image"]); image != "" {
		metadata["image"] = image
	}
	if source.Name != "" {
		metadata["sourceName"] = source.Name
	}

	sanitizedSummary := sanitizeString(s.sanitizer, summary)
	sanitizedContent := sanitizeString(s.sanitizer, content)

	hash := computeHash(pageURL.String(), title, sanitizedSummary, sanitizedContent)

	kind := inferKind(pageURL.Path, typeVal)

	item := models.ExternalItem{
		ID:       uuid.Nil,
		SourceID: source.ID,
		Kind:     kind,
		Title:    title,
		URL:      pageURL.String(),
		Metadata: metadata,
		Hash:     hash,
		Visible:  true,
	}
	if sanitizedSummary != "" {
		item.Summary = &sanitizedSummary
	}
	if sanitizedContent != "" {
		item.Content = &sanitizedContent
	}
	if publishedAt != nil {
		item.PublishedAt = publishedAt
	}

	return &item
}

func (s *Service) fetch(ctx context.Context, target *url.URL, cond *conditionalHeaders, skipRobot bool) (*http.Response, []byte, *conditionalHeaders, error) {
	if target == nil {
		return nil, nil, nil, errors.New("target url nil")
	}
	if !skipRobot {
		if err := s.ensureAllowed(ctx, target); err != nil {
			return nil, nil, nil, err
		}
	}

	var lastErr error
	for attempt := 0; attempt < s.retryCount; attempt++ {
		if err := s.limiter.Wait(ctx); err != nil {
			return nil, nil, nil, err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.String(), nil)
		if err != nil {
			return nil, nil, nil, err
		}
		req.Header.Set("User-Agent", s.userAgent)
		if cond != nil {
			if cond.etag != nil {
				req.Header.Set("If-None-Match", *cond.etag)
			}
			if cond.lastModified != nil {
				req.Header.Set("If-Modified-Since", cond.lastModified.UTC().Format(http.TimeFormat))
			}
		}

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, err := readBody(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		resp.Body = io.NopCloser(bytes.NewReader(body))

		meta := &conditionalHeaders{}
		if etag := resp.Header.Get("ETag"); etag != "" {
			meta.etag = &etag
		}
		if lm := resp.Header.Get("Last-Modified"); lm != "" {
			if t, err := http.ParseTime(lm); err == nil {
				meta.lastModified = &t
			}
		}

		return resp, body, meta, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("failed to fetch %s", target)
	}
	return nil, nil, nil, lastErr
}

func (s *Service) ensureAllowed(ctx context.Context, target *url.URL) error {
	if err := s.ensureHostAllowed(target.Host); err != nil {
		return err
	}
	return s.ensureRobots(ctx, target)
}

func (s *Service) ensureHostAllowed(host string) error {
	host = strings.ToLower(host)
	if host == "" {
		return errors.New("empty host")
	}
	for allowed := range s.allowlist {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return nil
		}
	}
	return fmt.Errorf("host %s not in allowlist", host)
}

func (s *Service) ensureRobots(ctx context.Context, target *url.URL) error {
	s.robotsMu.Lock()
	data, ok := s.robots[target.Host]
	s.robotsMu.Unlock()

	if !ok {
		robotsURL := &url.URL{Scheme: target.Scheme, Host: target.Host, Path: "/robots.txt"}
		resp, body, _, err := s.fetch(ctx, robotsURL, nil, true)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			data = nil
		} else {
			parsed, err := robotstxt.FromStatusAndBytes(resp.StatusCode, body)
			if err != nil {
				return err
			}
			data = parsed
		}
		s.robotsMu.Lock()
		s.robots[target.Host] = data
		s.robotsMu.Unlock()
	}

	if data == nil {
		return nil
	}
	group := data.FindGroup(s.userAgent)
	if group == nil {
		group = data.FindGroup("*")
	}
	if group == nil {
		return nil
	}
	if !group.Test(target.Path) {
		return fmt.Errorf("robots disallow %s", target)
	}
	return nil
}

func sanitizeString(policy *bluemonday.Policy, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	cleaned := policy.Sanitize(value)
	return strings.TrimSpace(cleaned)
}

func computeHash(values ...string) string {
	hasher := sha256.New()
	for _, value := range values {
		hasher.Write([]byte(value))
		hasher.Write([]byte{0})
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func stringFrom(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case json.Number:
		return v.String()
	default:
		return ""
	}
}

func parseTime(raw string) *time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	layouts := []string{time.RFC3339, "2006-01-02", time.RFC1123, "2006-01"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return &t
		}
	}
	return nil
}

func inferKind(path, typeVal string) string {
	lower := strings.ToLower(path)
	switch {
	case strings.Contains(lower, "/project"):
		return "project"
	case strings.Contains(lower, "/service"):
		return "service"
	case strings.Contains(lower, "/blog") || strings.Contains(lower, "/post"):
		return "post"
	}

	switch strings.ToLower(typeVal) {
	case "creativework", "project", "portfolio":
		return "project"
	case "service", "offer":
		return "service"
	default:
		return "post"
	}
}

func readBody(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, errors.New("nil body")
	}
	limited := io.LimitReader(body, maxBodyBytes)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	return data, nil
}
