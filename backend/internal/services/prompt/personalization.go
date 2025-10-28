package prompt

import (
	"fmt"
	"math"
	"strings"

	"github.com/tanydotai/tanyai/backend/internal/embedding"
	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

// BuildPersonalizedPrompt augments the grounded prompt with personalization snippets when available.
func BuildPersonalizedPrompt(base kb.KnowledgeBase, question string, personalization embedding.PersonalizationResult) string {
	prompt := BuildPrompt(base, question)
	if !personalization.Enabled || personalization.Weight <= 0 || len(personalization.Snippets) == 0 {
		return prompt
	}

	weight := personalization.Weight
	percent := int(math.Round(weight * 100))
	if percent <= 0 {
		percent = 1
	}

	var builder strings.Builder
	builder.WriteString(prompt)
	builder.WriteString("\n\nInstruksi personalisasi:")
	builder.WriteString("\n- Terapkan gaya dan nada khas dengan bobot sekitar ")
	builder.WriteString(fmt.Sprintf("%d%%.\n", percent))

	for _, snippet := range personalization.Snippets {
		builder.WriteString(formatSnippet(snippet))
		builder.WriteString("\n")
	}
	builder.WriteString("Pastikan jawaban tetap akurat dan profesional sambil menjaga karakter personal.")
	return builder.String()
}

func formatSnippet(snippet embedding.Snippet) string {
	kind := strings.ToUpper(snippet.Kind)
	if kind == "" {
		kind = "KONTEKS"
	}
	score := snippet.Score
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	return fmt.Sprintf("- [%s • skor %.2f] %s", kind, score, strings.TrimSpace(snippet.Content))
}
