package prompt

import (
        "strings"
        "testing"
        "time"

        "github.com/tanydotai/tanyai/backend/internal/services/kb"
)

func sampleBase() kb.KnowledgeBase {
        return kb.KnowledgeBase{
                Profile: kb.Profile{Name: "Tanya", Title: "Freelance Engineer", Bio: "Bio", Location: "Jakarta", Email: "tanya@example.com"},
                Services: []kb.Service{
                        {Name: "Build", Description: "Build apps", PriceRange: []string{"10jt"}, DurationLabel: "2 minggu", Order: 1},
			{Name: "Optimize", Description: "Optimize systems", Order: 2},
			{Name: "Consult", Description: "Consulting", Order: 3},
			{Name: "Extra", Description: "Extra", Order: 4},
		},
                Projects: []kb.Project{
                        {Title: "Project A", Description: "Desc", TechStack: []string{"Go", "React"}, IsFeatured: true, Order: 1},
                        {Title: "Project B", Description: "Another", TechStack: []string{"Next.js"}, Order: 2},
                },
                Posts: []kb.Post{
                        {Title: "New Launch", Summary: "Diluncurkan fitur baru", URL: "https://example.com/post", Source: "noahis.me", PublishedAt: time.Now()},
                },
        }
}

func TestBuildPromptIncludesContextAndQuestion(t *testing.T) {
	prompt := BuildPrompt(sampleBase(), "Apa layananmu?")
	if !strings.Contains(prompt, "Tanya") {
		t.Fatalf("prompt should include profile name")
	}
	if strings.Count(prompt, "- ") < 3 {
		t.Fatalf("prompt should list services")
	}
	if !strings.Contains(prompt, "Pertanyaan: Apa layananmu?") {
		t.Fatalf("prompt should include question section")
	}
        if strings.Contains(prompt, "Extra") {
                t.Fatalf("prompt should limit number of services")
        }
        if !strings.Contains(prompt, "Update terbaru") {
                t.Fatalf("prompt should include external updates")
        }
}

func TestSummarizeForHumanReferencesFeaturedProject(t *testing.T) {
        response := SummarizeForHuman("Apa layananmu?", sampleBase())
        if !strings.Contains(response, "Project A") {
                t.Fatalf("expected featured project to be referenced")
        }
        if !strings.Contains(response, "Build") {
                t.Fatalf("expected service mention in summary")
        }
        if !strings.Contains(response, "New Launch") {
                t.Fatalf("expected latest post mention in summary")
        }
}
