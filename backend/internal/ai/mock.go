package ai

import "context"

// Mock implements Provider and returns a deterministic response for testing.
type Mock struct {
	Text string
}

// NewMock creates a mock provider returning a canned response.
func NewMock() *Mock {
	return &Mock{Text: "Mock response"}
}

// Generate returns the configured mock text regardless of the request payload.
func (m *Mock) Generate(_ context.Context, _ Request) (Response, error) {
	text := m.Text
	if text == "" {
		text = "Mock response"
	}
	return Response{Text: text}, nil
}
