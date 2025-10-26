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
	defaultGeminiEndpoint = "https://generativelanguage.googleapis.com"
	defaultGeminiModel    = "gemini-2.5-pro"
)

// Gemini implements the Provider interface using Google Gemini's REST API.
type Gemini struct {
	Key      string
	Model    string
	Client   *http.Client
	Endpoint string
}

// NewGemini constructs a Gemini provider with the supplied API key and model.
func NewGemini(key, model string) *Gemini {
	model = strings.TrimSpace(model)
	if model == "" {
		model = defaultGeminiModel
	}

	return &Gemini{
		Key:      strings.TrimSpace(key),
		Model:    model,
		Client:   &http.Client{Timeout: 15 * time.Second},
		Endpoint: defaultGeminiEndpoint,
	}
}

// Generate invokes the Gemini API to create a text response.
func (g *Gemini) Generate(ctx context.Context, r Request) (Response, error) {
	if strings.TrimSpace(r.Prompt) == "" {
		return Response{}, errors.New("prompt is required")
	}
	if g == nil {
		return Response{}, errors.New("gemini provider is not configured")
	}
	if strings.TrimSpace(g.Key) == "" {
		return Response{}, errors.New("missing GOOGLE_GENAI_API_KEY")
	}

	maxTokens := r.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 512
	}
	temperature := r.Temperature
	if temperature <= 0 {
		temperature = 0.4
	}

	payload := map[string]any{
		"contents": []any{
			map[string]any{
				"parts": []any{
					map[string]any{"text": r.Prompt},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":     temperature,
			"maxOutputTokens": maxTokens,
			"candidateCount":  1,
			"topP":           0.8,
			"topK":           20,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, fmt.Errorf("failed to encode request: %w", err)
	}

	endpoint := g.Endpoint
	if endpoint == "" {
		endpoint = defaultGeminiEndpoint
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

	// Debug output
	fmt.Printf("Gemini request to: %s/v1beta/models/%s:generateContent\n", endpoint, g.Model)
	fmt.Printf("Request body: %s\n", string(body))

	apiURL := fmt.Sprintf("%s/v1beta/models/%s:generateContent", endpoint, g.Model)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return Response{}, fmt.Errorf("failed to create request: %w", err)
	}

	query := req.URL.Query()
	query.Set("key", g.Key)
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Type", "application/json")

	client := g.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return Response{}, fmt.Errorf(
			"gemini request failed: status=%d body=%s",
			resp.StatusCode,
			strings.TrimSpace(string(data)),
		)
	}

	// Read and log the response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to read response body: %w", err)
	}
	
	fmt.Printf("Gemini API Response: %s\n", string(bodyBytes))
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var decoded struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"error"`
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			FinishReason string `json:"finishReason"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&decoded); err != nil {
		return Response{}, fmt.Errorf("failed to decode response: %w\nResponse body: %s", err, string(bodyBytes))
	}

	if decoded.Error.Message != "" {
		return Response{}, fmt.Errorf("gemini API error: %s", decoded.Error.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return Response{}, fmt.Errorf("failed to decode gemini response: %w", err)
	}

	if len(decoded.Candidates) == 0 || len(decoded.Candidates[0].Content.Parts) == 0 {
		return Response{}, errors.New("empty response from gemini")
	}

	text := strings.TrimSpace(decoded.Candidates[0].Content.Parts[0].Text)
	if text == "" {
		return Response{}, errors.New("empty response from gemini")
	}

	return Response{Text: text}, nil
}
