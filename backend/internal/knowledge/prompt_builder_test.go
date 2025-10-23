package knowledge

import "strings"
import "testing"

func TestBuildSystemPromptContainsCoreData(t *testing.T) {
	base := LoadStaticKnowledgeBase()
	prompt := BuildSystemPrompt(base)

	if !strings.Contains(prompt, base.Profile.Name) {
		t.Fatalf("prompt should mention profile name")
	}
	if len(base.Services) == 0 {
		t.Fatalf("fixtures should include services")
	}
	if !strings.Contains(prompt, base.Services[0].Name) {
		t.Fatalf("prompt should list service names")
	}
	if !strings.Contains(prompt, base.Contact.Email) {
		t.Fatalf("prompt should include contact info")
	}
}

func TestSummarizeForHumanReferencesQuestion(t *testing.T) {
	base := LoadStaticKnowledgeBase()
	summary := SummarizeForHuman("Apa layananmu?", base)

	if !strings.Contains(summary, "Apa layananmu?") {
		t.Fatalf("summary should contain the original question")
	}
	if !strings.Contains(summary, base.Profile.Name) {
		t.Fatalf("summary should introduce the profile")
	}
}
