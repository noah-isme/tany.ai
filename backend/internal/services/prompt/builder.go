package prompt

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

// BuildPrompt creates a grounded single-message prompt suitable for Gemini style inputs.
func BuildPrompt(base kb.KnowledgeBase, question string) string {
	var builder strings.Builder
	profile := base.Profile

	builder.WriteString("Anda adalah asisten virtual untuk ")
	if profile.Name != "" {
		builder.WriteString(profile.Name)
	} else {
		builder.WriteString("tany.ai")
	}
	builder.WriteString(". Jawab menggunakan informasi berikut.\n\n")

	if profile.Title != "" || profile.Location != "" || profile.Bio != "" {
		builder.WriteString("Profil singkat:\n")
		if profile.Title != "" {
			builder.WriteString(fmt.Sprintf("- Peran: %s\n", profile.Title))
		}
		if profile.Location != "" {
			builder.WriteString(fmt.Sprintf("- Lokasi: %s\n", profile.Location))
		}
		if profile.Bio != "" {
			builder.WriteString(fmt.Sprintf("- Bio: %s\n", profile.Bio))
		}
		builder.WriteString("\n")
	}

	services := topServices(base.Services, 3)
	if len(services) > 0 {
		builder.WriteString("Layanan prioritas:\n")
		for _, service := range services {
			line := fmt.Sprintf("- %s", service.Name)
			details := make([]string, 0, 3)
			if service.Description != "" {
				details = append(details, service.Description)
			}
			if len(service.PriceRange) > 0 {
				currency := service.Currency
				if currency == "" {
					currency = "IDR"
				}
				details = append(details, fmt.Sprintf("Harga %s %s", currency, strings.Join(service.PriceRange, " – ")))
			}
			if service.DurationLabel != "" {
				details = append(details, fmt.Sprintf("Durasi %s", service.DurationLabel))
			}
			if len(details) > 0 {
				line += " — " + strings.Join(details, "; ")
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	projects := topProjects(base.Projects, 2)
	if len(projects) > 0 {
		builder.WriteString("Portofolio unggulan:\n")
		for _, project := range projects {
			line := fmt.Sprintf("- %s", project.Title)
			details := make([]string, 0, 3)
			if project.Description != "" {
				details = append(details, project.Description)
			}
			if len(project.TechStack) > 0 {
				details = append(details, fmt.Sprintf("Tech: %s", strings.Join(project.TechStack, ", ")))
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

	builder.WriteString("Aturan:\n")
	builder.WriteString("1. Gunakan hanya informasi di atas.\n")
	builder.WriteString("2. Jika data tidak tersedia, jawab bahwa informasinya belum ada.\n")
	builder.WriteString("3. Jawab dalam bahasa Indonesia yang ramah profesional dan ringkas.\n\n")

	builder.WriteString("Pertanyaan: ")
	builder.WriteString(strings.TrimSpace(question))
	builder.WriteString("\n\nJawaban yang relevan:")

	return builder.String()
}

func topServices(services []kb.Service, limit int) []kb.Service {
	if len(services) == 0 {
		return nil
	}
	sorted := make([]kb.Service, len(services))
	copy(sorted, services)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Order < sorted[j].Order
	})
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func topProjects(projects []kb.Project, limit int) []kb.Project {
	if len(projects) == 0 {
		return nil
	}
	sorted := make([]kb.Project, len(projects))
	copy(sorted, projects)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].IsFeatured == sorted[j].IsFeatured {
			return sorted[i].Order < sorted[j].Order
		}
		return sorted[i].IsFeatured && !sorted[j].IsFeatured
	})
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

// SummarizeForHuman returns a deterministic answer if the provider is unavailable.
func SummarizeForHuman(question string, base kb.KnowledgeBase) string {
	services := topServices(base.Services, 3)
	serviceNames := make([]string, 0, len(services))
	for _, service := range services {
		serviceNames = append(serviceNames, service.Name)
	}

	projects := topProjects(base.Projects, 1)
	featured := ""
	if len(projects) > 0 {
		featured = projects[0].Title
	}

	contact := ""
	if base.Profile.Email != "" {
		contact = fmt.Sprintf("Hubungi %s", base.Profile.Email)
	} else if base.Profile.Phone != "" {
		contact = fmt.Sprintf("Kontak %s", base.Profile.Phone)
	}

	builder := strings.Builder{}
	builder.WriteString("Pertanyaan diterima: ")
	builder.WriteString(strings.TrimSpace(question))
	builder.WriteString("\n\n")
	if base.Profile.Name != "" {
		builder.WriteString(fmt.Sprintf("Halo! Saya %s. ", base.Profile.Name))
	} else {
		builder.WriteString("Halo! Saya asisten tany.ai. ")
	}
	if len(serviceNames) > 0 {
		builder.WriteString("Saat ini saya menawarkan ")
		builder.WriteString(strings.Join(serviceNames, ", "))
		builder.WriteString(". ")
	}
	if featured != "" {
		builder.WriteString(fmt.Sprintf("Contoh proyek terbaru: %s. ", featured))
	}
	if contact != "" {
		builder.WriteString(contact)
		builder.WriteString(" untuk detail lanjut.")
	}

	return strings.TrimSpace(builder.String())
}
