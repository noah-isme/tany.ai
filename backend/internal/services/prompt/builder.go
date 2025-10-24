package prompt

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

// BuildSystemPrompt constructs a grounded system prompt for the chat model.
func BuildSystemPrompt(base kb.KnowledgeBase) string {
	var builder strings.Builder
	profile := base.Profile
	builder.WriteString("Kamu adalah asisten virtual untuk ")
	builder.WriteString(profile.Name)
	builder.WriteString(". Jawab menggunakan data internal berikut dan jangan berimprovisasi.\n\n")

	if profile.Title != "" || profile.Bio != "" || profile.Location != "" {
		builder.WriteString("Profil:\n")
		if profile.Title != "" {
			builder.WriteString(fmt.Sprintf("- Title: %s\n", profile.Title))
		}
		if profile.Bio != "" {
			builder.WriteString(fmt.Sprintf("- Bio: %s\n", profile.Bio))
		}
		if profile.Location != "" {
			builder.WriteString(fmt.Sprintf("- Lokasi: %s\n", profile.Location))
		}
		builder.WriteString("\n")
	}

	if len(base.Skills) > 0 {
		skills := make([]string, 0, len(base.Skills))
		for _, skill := range base.Skills {
			skills = append(skills, skill.Name)
		}
		builder.WriteString("Keahlian utama: ")
		builder.WriteString(strings.Join(skills, ", "))
		builder.WriteString(".\n\n")
	}

	if len(base.Services) > 0 {
		builder.WriteString("Layanan yang tersedia:\n")
		for _, service := range base.Services {
			line := fmt.Sprintf("- %s", service.Name)
			details := make([]string, 0, 3)
			if service.Description != "" {
				details = append(details, service.Description)
			}
			if len(service.PriceRange) > 0 {
				details = append(details, fmt.Sprintf("Harga: %s", strings.Join(service.PriceRange, " - ")))
			}
			if service.DurationLabel != "" {
				details = append(details, fmt.Sprintf("Durasi: %s", service.DurationLabel))
			}
			if len(details) > 0 {
				line += " — " + strings.Join(details, "; ")
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	if len(base.Projects) > 0 {
		builder.WriteString("Proyek unggulan:\n")
		projects := make([]kb.Project, len(base.Projects))
		copy(projects, base.Projects)
		sort.SliceStable(projects, func(i, j int) bool {
			if projects[i].IsFeatured == projects[j].IsFeatured {
				return projects[i].Order < projects[j].Order
			}
			return projects[i].IsFeatured && !projects[j].IsFeatured
		})
		for _, project := range projects {
			line := fmt.Sprintf("- %s", project.Title)
			details := make([]string, 0, 3)
			if project.Description != "" {
				details = append(details, project.Description)
			}
			if len(project.TechStack) > 0 {
				details = append(details, fmt.Sprintf("Teknologi: %s", strings.Join(project.TechStack, ", ")))
			}
			if project.ProjectURL != "" {
				details = append(details, fmt.Sprintf("URL: %s", project.ProjectURL))
			}
			if len(details) > 0 {
				line += " — " + strings.Join(details, "; ")
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("Aturan ketat:\n")
	builder.WriteString("1. Gunakan hanya informasi di atas.\n")
	builder.WriteString("2. Jika data tidak tersedia, jawab dengan jujur bahwa informasinya belum ada.\n")
	builder.WriteString("3. Jawab dalam bahasa Indonesia yang ramah profesional.\n")

	return builder.String()
}

// SummarizeForHuman provides a deterministic fallback answer using the knowledge base.
func SummarizeForHuman(question string, base kb.KnowledgeBase) string {
	var services []string
	for _, service := range base.Services {
		services = append(services, service.Name)
	}
	var featured string
	for _, project := range base.Projects {
		if project.IsFeatured {
			featured = project.Title
			break
		}
	}
	if featured == "" && len(base.Projects) > 0 {
		featured = base.Projects[0].Title
	}

	return fmt.Sprintf("Pertanyaan diterima: %s\n\nHalo! Saya %s. Saya menawarkan layanan %s. Contoh proyek terbaru: %s.\nSilakan hubungi kami untuk detail lanjutan.",
		question,
		base.Profile.Name,
		strings.Join(services, ", "),
		featured,
	)
}
