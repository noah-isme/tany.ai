package prompt

import (
	"strings"
	"testing"

	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

func sampleBase() kb.KnowledgeBase {
	return kb.KnowledgeBase{
		Profile:  kb.Profile{Name: "Tanya", Title: "Freelance Engineer", Bio: "Bio", Location: "Jakarta"},
		Skills:   []kb.Skill{{Name: "Go"}, {Name: "Next.js"}},
		Services: []kb.Service{{Name: "Build", Description: "Build apps", PriceRange: []string{"IDR 10jt"}, DurationLabel: "2 minggu", Order: 1}},
		Projects: []kb.Project{{Title: "Project A", Description: "Desc", TechStack: []string{"Go", "React"}, IsFeatured: true, Order: 1}},
	}
}

func TestBuildSystemPromptIncludesProfile(t *testing.T) {
	prompt := BuildSystemPrompt(sampleBase())
	if !strings.Contains(prompt, "Tanya") {
		t.Fatalf("prompt should include profile name")
	}
	if !strings.Contains(prompt, "Keahlian utama") {
		t.Fatalf("prompt should include skills section")
	}
}

func TestSummarizeForHumanReferencesFeaturedProject(t *testing.T) {
	response := SummarizeForHuman("Apa layananmu?", sampleBase())
	if !strings.Contains(response, "Project A") {
		t.Fatalf("expected featured project to be referenced")
	}
}
