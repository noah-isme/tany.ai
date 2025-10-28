package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	defaultOpenAIEndpoint       = "https://api.openai.com/v1"
	defaultOpenAIEmbeddingModel = "text-embedding-3-large"
)

// OpenAIEmbedding implements embedding generation via the OpenAI REST API.
type OpenAIEmbedding struct {
	Key      string
	Model    string
	Endpoint string
	Client   *http.Client
}

// NewOpenAIEmbedding constructs an embedding client using the provided API key and model.
func NewOpenAIEmbedding(key, model string) *OpenAIEmbedding {
	model = strings.TrimSpace(model)
	if model == "" {
		model = defaultOpenAIEmbeddingModel
	}
	return &OpenAIEmbedding{
		Key:      strings.TrimSpace(key),
		Model:    model,
		Endpoint: defaultOpenAIEndpoint,
		Client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Embed creates a vector representation for the supplied text input.
func (o *OpenAIEmbedding) Embed(ctx context.Context, input string) ([]float32, error) {
	if o == nil {
		return nil, errors.New("openai embedding client is nil")
	}
	if strings.TrimSpace(o.Key) == "" {
		return nil, errors.New("OPENAI_API_KEY is required for embeddings")
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, errors.New("input text is required")
	}

	payload := map[string]any{
		"input": input,
		"model": o.Model,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode openai embedding request: %w", err)
	}

	endpoint := strings.TrimSuffix(o.Endpoint, "/")
	if endpoint == "" {
		endpoint = defaultOpenAIEndpoint
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create openai embedding request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.Key)

	client := o.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var data struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&data)
		msg := data.Error.Message
		if msg == "" {
			msg = resp.Status
		}
		return nil, fmt.Errorf("openai embeddings request failed: %s", msg)
	}

	var decoded struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode openai embedding response: %w", err)
	}
	if len(decoded.Data) == 0 || len(decoded.Data[0].Embedding) == 0 {
		return nil, errors.New("openai embedding response missing data")
	}

	floats := make([]float32, len(decoded.Data[0].Embedding))
	for i, v := range decoded.Data[0].Embedding {
		floats[i] = float32(v)
	}
	return floats, nil
}
