package ai

import "context"

// Request captures a text generation request sent to an AI provider.
type Request struct {
	Prompt      string
	MaxTokens   int
	Temperature float32
}

// Response represents a normalized AI generation result.
type Response struct {
	Text string
}

// Provider describes the capabilities required from any AI text generator.
type Provider interface {
	Generate(ctx context.Context, r Request) (Response, error)
}
