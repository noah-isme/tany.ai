package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultLeapcellEndpoint = "https://api.leapcell.io/v1"
)

// Leapcell implements the Provider interface using Leapcell's API.
type Leapcell struct {
	Key      string
	Client   *http.Client
	Endpoint string
	ProjectID string
	TableID   string
}

// NewLeapcell constructs a Leapcell provider with the supplied API key.
func NewLeapcell(key, projectID, tableID string) *Leapcell {
	return &Leapcell{
		Key:       strings.TrimSpace(key),
		ProjectID: strings.TrimSpace(projectID),
		TableID:   strings.TrimSpace(tableID),
		Client:    &http.Client{Timeout: 15 * time.Second},
		Endpoint:  defaultLeapcellEndpoint,
	}
}

// Generate invokes the Leapcell API to create a text response.
func (l *Leapcell) Generate(ctx context.Context, r Request) (Response, error) {
	if strings.TrimSpace(r.Prompt) == "" {
		return Response{}, errors.New("prompt is required")
	}
	if l == nil {
		return Response{}, errors.New("leapcell provider is not configured")
	}
	if strings.TrimSpace(l.Key) == "" {
		return Response{}, errors.New("missing LEAPCELL_API_KEY")
	}
	if l.ProjectID == "" || l.TableID == "" {
		return Response{}, errors.New("missing project_id or table_id")
	}

	payload := map[string]any{
		"question": r.Prompt,
		"project_id": l.ProjectID,
		"table_id": l.TableID,
		"config": map[string]any{
			"temperature": r.Temperature,
			"max_tokens":  r.MaxTokens,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("failed to encode request: %w", err)
	}

	endpoint := l.Endpoint
	if endpoint == "" {
		endpoint = defaultLeapcellEndpoint
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/chat", endpoint),
		bytes.NewReader(body),
	)
	if err != nil {
		return Response{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.Key)

	resp, err := l.Client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Response{}, fmt.Errorf("leapcell API error: status=%d body=%s", resp.StatusCode, string(body))
	}

	var result struct {
		Answer string `json:"answer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Response{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return Response{Text: result.Answer}, nil
}