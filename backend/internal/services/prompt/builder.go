package prompt

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/tanydotai/tanyai/backend/internal/services/kb"
)

const (
	defaultMaxServicesInPrompt = 3
	defaultMaxProjectsInPrompt = 3
)

func maxFromEnv(key string, fallback int) int {
	if fallback <= 0 {
		fallback = 1
	}
	if raw := os.Getenv(key); raw != "" {
		if val, err := strconv.Atoi(raw); err == nil && val > 0 {
			return val
		}
	}
	return fallback
}

// BuildPrompt creates a grounded single-message prompt suitable for Gemini style inputs.
// Empty or whitespace questions will return an error about invalid input.
func BuildPrompt(base kb.KnowledgeBase, question string) string {
	// Validate and sanitize input
	question = strings.TrimSpace(question)
	if question == "" {
		return "Mohon maaf, saya tidak dapat memproses pertanyaan kosong. Silakan ajukan pertanyaan Anda."
	}

	questionLower := strings.ToLower(question)

	// Handle special case for new users
	if questionLower == "hi" || questionLower == "halo" || questionLower == "hello" {
		question = "Perkenalkan diri Anda dan layanan yang tersedia"
		questionLower = strings.ToLower(question)
	}
	var builder strings.Builder
	profile := base.Profile

	var maxLen int
	if strings.Contains(questionLower, "layanan") {
		maxLen = 800 // Lebih banyak ruang untuk informasi layanan
	} else if strings.Contains(questionLower, "proyek") || strings.Contains(questionLower, "portfolio") {
		maxLen = 600 // Lebih banyak ruang untuk informasi proyek
	} else {
		maxLen = 400 // Default untuk pertanyaan umum
	}

	// Start with the question
	builder.WriteString("Pertanyaan: ")
	if len(question) > 100 {
		builder.WriteString(question[:97] + "...")
	} else {
		builder.WriteString(question)
	}
	builder.WriteString("\n\n")

	// Add context header
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
			bio := profile.Bio
			if len(bio) > 200 {
				bio = bio[:197] + "..."
			}
			builder.WriteString(fmt.Sprintf("- Bio: %s\n", bio))
		}
		builder.WriteString("\n")
	}

	maxServicesAllowed := maxFromEnv("PROMPT_MAX_SERVICES", defaultMaxServicesInPrompt)
	serviceLimit := maxServicesAllowed
	if serviceLimit <= 0 {
		serviceLimit = defaultMaxServicesInPrompt
	}
	serviceFocused := strings.Contains(questionLower, "layanan") ||
		strings.Contains(questionLower, "jasa") ||
		strings.Contains(questionLower, "service")
	if !serviceFocused && serviceLimit > 1 {
		serviceLimit = 1
	}
	services := topServices(base.Services, serviceLimit)
	if len(services) > 0 {
		builder.WriteString("Layanan prioritas:\n")
		for _, service := range services {
			line := fmt.Sprintf("- %s", service.Name)
			details := make([]string, 0, 3)
			if service.Description != "" {
				desc := service.Description
				if len(desc) > 100 {
					desc = desc[:97] + "..."
				}
				details = append(details, desc)
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
			if len(line) > 200 {
				line = line[:197] + "..."
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	maxProjectsAllowed := maxFromEnv("PROMPT_MAX_PROJECTS", defaultMaxProjectsInPrompt)
	if maxProjectsAllowed <= 0 {
		maxProjectsAllowed = defaultMaxProjectsInPrompt
	}
	projectFocused := strings.Contains(questionLower, "proyek") || strings.Contains(questionLower, "portfolio")
	projectLimit := 1
	if projectFocused {
		projectLimit = maxProjectsAllowed
	} else if maxProjectsAllowed < projectLimit {
		projectLimit = maxProjectsAllowed
	}
	projects := topProjects(base.Projects, projectLimit)
	if len(projects) > 0 {
		builder.WriteString("Portofolio unggulan:\n")
		for _, project := range projects {
			line := fmt.Sprintf("- %s", project.Title)
			details := make([]string, 0, 3)
			if project.Description != "" {
				desc := project.Description
				if len(desc) > 100 {
					desc = desc[:97] + "..."
				}
				details = append(details, desc)
			}
			if len(project.TechStack) > 0 {
				stack := project.TechStack
				if len(stack) > 5 {
					stack = stack[:5]
				}
				details = append(details, fmt.Sprintf("Tech: %s", strings.Join(stack, ", ")))
			}
			if project.ProjectURL != "" {
				details = append(details, fmt.Sprintf("URL: %s", project.ProjectURL))
			}
			if len(details) > 0 {
				line += " — " + strings.Join(details, "; ")
			}
			if len(line) > 200 {
				line = line[:197] + "..."
			}
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	// Add final instructions
	builder.WriteString("\nInstruksi: Jawab dengan ringkas dan ramah dalam bahasa Indonesia. Gunakan hanya informasi yang tersedia di atas.\n\n")
	builder.WriteString("Berikan jawaban untuk: ")
	builder.WriteString(strings.TrimSpace(question))

	// Get final text and ensure it's complete
	result := builder.String()

	// If we need to trim, do it intelligently to keep important parts
	if len(result) > maxLen {
		const summaryFormat = `Anda adalah %s, %s yang berlokasi di %s.

%s

%s

Berikan jawaban untuk: %s`

		// Get core information
		role := "Full Stack Developer & AI Consultant"
		if profile.Title != "" {
			role = profile.Title
		}

		location := "Indonesia"
		if profile.Location != "" {
			location = profile.Location
		}

		// Get top service
		serviceDesc := "Tidak ada informasi layanan"
		if len(services) > 0 {
			details := make([]string, 0, 2)

			if len(services[0].PriceRange) > 0 {
				currency := services[0].Currency
				if currency == "" {
					currency = "IDR"
				}
				priceRange := strings.Join(services[0].PriceRange, " - ")
				priceRange = strings.ReplaceAll(priceRange, "IDR ", "") // Remove duplicate IDR
				details = append(details, fmt.Sprintf("Harga %s %s", currency, priceRange))
			}

			if services[0].DurationLabel != "" {
				details = append(details, fmt.Sprintf("Durasi %s", services[0].DurationLabel))
			}

			detailStr := ""
			if len(details) > 0 {
				detailStr = fmt.Sprintf(" (%s)", strings.Join(details, ", "))
			}

			serviceDesc = fmt.Sprintf("Layanan utama:\n- %s%s",
				services[0].Name,
				detailStr)
		}

		// Format the shortened version
		return fmt.Sprintf(summaryFormat,
			profile.Name,
			role,
			location,
			serviceDesc,
			"Instruksi: Jawab dengan ringkas dan ramah dalam bahasa Indonesia. Gunakan hanya informasi yang tersedia di atas.",
			strings.TrimSpace(question))
	}

	return result
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
