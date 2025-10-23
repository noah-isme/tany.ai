package knowledge

import (
	"fmt"
	"strings"
)

// BuildSystemPrompt constructs the base system prompt used to ground the LLM
// with freelancer specific knowledge.
func BuildSystemPrompt(base KnowledgeBase) string {
	var builder strings.Builder
	builder.WriteString("Kamu adalah asisten virtual untuk ")
	builder.WriteString(base.Profile.Name)
	builder.WriteString(". Gunakan fakta berikut untuk menjawab pertanyaan klien secara profesional.\n\n")

	builder.WriteString("Profil:\n")
	builder.WriteString(fmt.Sprintf("- Tagline: %s\n", base.Profile.Tagline))
	builder.WriteString(fmt.Sprintf("- Bio: %s\n", base.Profile.Bio))
	builder.WriteString(fmt.Sprintf("- Keahlian: %s\n", strings.Join(base.Profile.Expertise, ", ")))
	builder.WriteString(fmt.Sprintf("- Berpengalaman: %d tahun\n", base.Profile.YearsActive))
	builder.WriteString(fmt.Sprintf("- Berbasis di: %s\n\n", base.Profile.Location))

	builder.WriteString("Layanan:\n")
	for _, service := range base.Services {
		builder.WriteString(fmt.Sprintf("- %s (mulai %s): %s. Deliverables: %s.\n",
			service.Name,
			service.StartingAt,
			service.Description,
			strings.Join(service.Deliverables, ", "),
		))
	}
	builder.WriteString("\nPortfolio Utama:\n")
	for _, project := range base.Projects {
		builder.WriteString(fmt.Sprintf("- %s menggunakan %s. Dampak: %s.\n",
			project.Title,
			strings.Join(project.TechStack, ", "),
			project.Impact,
		))
	}

	builder.WriteString("\nStruktur Harga:\n")
	for _, tier := range base.Pricing {
		builder.WriteString(fmt.Sprintf("- Paket %s (%s): %s.\n",
			tier.Name,
			tier.Price,
			strings.Join(tier.WhatYouGet, ", "),
		))
	}

	builder.WriteString("\nKontak:\n")
	builder.WriteString(fmt.Sprintf("- Email: %s\n", base.Contact.Email))
	builder.WriteString(fmt.Sprintf("- Website: %s\n", base.Contact.Website))
	builder.WriteString(fmt.Sprintf("- WhatsApp: %s\n", base.Contact.WhatsApp))
	builder.WriteString(fmt.Sprintf("- LinkedIn: %s", base.Contact.LinkedIn))

	builder.WriteString("\n\nAturan:\n")
	builder.WriteString("1. Hanya gunakan informasi di atas.\n")
	builder.WriteString("2. Jawab dengan nada ramah, profesional, dan ringkas.\n")
	builder.WriteString("3. Jika pertanyaan di luar scope, arahkan klien untuk menghubungi via email.")

	return builder.String()
}

// SummarizeForHuman composes a friendly summary that dapat dikirim langsung ke user
// tanpa memanggil API OpenAI. Berguna untuk mock response selama pengembangan awal.
func SummarizeForHuman(question string, base KnowledgeBase) string {
	return fmt.Sprintf("Pertanyaan diterima: %s\n\nHalo! Saya %s. Saya menawarkan layanan %s. \nContoh proyek terbaru: %s.\nUntuk informasi lengkap atau diskusi lanjut, silakan hubungi saya di %s.",
		question,
		base.Profile.Name,
		joinNames(base.Services),
		base.Projects[0].Title,
		base.Contact.Email,
	)
}

func joinNames(services []Service) string {
	names := make([]string, 0, len(services))
	for _, service := range services {
		names = append(names, service.Name)
	}
	return strings.Join(names, ", ")
}
