package knowledge

// Profile contains the public profile information for the freelancer.
type Profile struct {
	Name        string   `json:"name"`
	Tagline     string   `json:"tagline"`
	Bio         string   `json:"bio"`
	Expertise   []string `json:"expertise"`
	YearsActive int      `json:"yearsActive"`
	Location    string   `json:"location"`
}

// Service represents a single offering that can be surfaced in responses.
type Service struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Deliverables []string `json:"deliverables"`
	Turnaround   string   `json:"turnaround"`
	StartingAt   string   `json:"startingAt"`
}

// Project highlights past work to give credibility in chat answers.
type Project struct {
	Title     string   `json:"title"`
	Summary   string   `json:"summary"`
	TechStack []string `json:"techStack"`
	Impact    string   `json:"impact"`
}

// PricingTier lists high level pricing information per package.
type PricingTier struct {
	Name       string   `json:"name"`
	Price      string   `json:"price"`
	WhatYouGet []string `json:"whatYouGet"`
}

// ContactInfo contains direct channels for leads.
type ContactInfo struct {
	Email    string `json:"email"`
	Website  string `json:"website"`
	WhatsApp string `json:"whatsApp"`
	LinkedIn string `json:"linkedIn"`
}

// KnowledgeBase aggregates all data powering the chat assistant.
type KnowledgeBase struct {
	Profile  Profile       `json:"profile"`
	Services []Service     `json:"services"`
	Projects []Project     `json:"projects"`
	Pricing  []PricingTier `json:"pricing"`
	Contact  ContactInfo   `json:"contact"`
}

// LoadStaticKnowledgeBase returns seed data that represents the freelancer.
// In production this should be replaced by persistence (database) integration.
func LoadStaticKnowledgeBase() KnowledgeBase {
	return KnowledgeBase{
		Profile: Profile{
			Name:        "Tanya A.I.",
			Tagline:     "AI Client Assistant for Modern Freelancers",
			Bio:         "Saya membantu menjawab pertanyaan calon klien secara profesional dengan konteks data freelancer.",
			Expertise:   []string{"Next.js", "Golang", "AI Prompt Engineering", "UI/UX Strategy"},
			YearsActive: 7,
			Location:    "Jakarta, Indonesia",
		},
		Services: []Service{
			{
				Name:         "Pembuatan Website Portfolio",
				Description:  "Website responsif dengan fokus pada storytelling brand personal.",
				Deliverables: []string{"Landing page", "Halaman layanan", "Integrasi formulir leads"},
				Turnaround:   "2-3 minggu",
				StartingAt:   "IDR 12.000.000",
			},
			{
				Name:         "AI Chat Assistant Custom",
				Description:  "Implementasi chatbot berbasis profil profesional untuk otomatisasi komunikasi klien.",
				Deliverables: []string{"Setup knowledge base", "Integrasi OpenAI GPT", "Dashboard monitoring"},
				Turnaround:   "3-4 minggu",
				StartingAt:   "IDR 18.000.000",
			},
			{
				Name:         "Consulting & Strategy",
				Description:  "Sesi konsultasi 1:1 mengenai positioning layanan digital dan optimasi workflow AI.",
				Deliverables: []string{"Audit profil", "Rencana peningkatan layanan", "Template komunikasi klien"},
				Turnaround:   "1 minggu",
				StartingAt:   "IDR 2.500.000",
			},
		},
		Projects: []Project{
			{
				Title:     "SaaS Landing Page untuk Startup Fintech",
				Summary:   "Membangun landing page interaktif dengan konversi tinggi untuk peluncuran produk fintech.",
				TechStack: []string{"Next.js", "Tailwind CSS", "Vercel"},
				Impact:    "Meningkatkan sign-up calon pengguna sebesar 32% dalam 3 minggu pertama.",
			},
			{
				Title:     "Chat Assistant untuk Agency Kreatif",
				Summary:   "Membuat chatbot khusus untuk menjawab pertanyaan layanan desain dan motion graphics.",
				TechStack: []string{"Golang", "MongoDB", "OpenAI GPT-4"},
				Impact:    "Mengurangi waktu respons tim sales dari 6 jam menjadi 15 menit rata-rata.",
			},
			{
				Title:     "Knowledge Hub Freelancer",
				Summary:   "Portal internal untuk mengelola portfolio, layanan, dan studi kasus.",
				TechStack: []string{"Next.js", "Supabase", "Clerk"},
				Impact:    "Mempercepat proses update portfolio bulanan dari 3 hari menjadi 1 hari kerja.",
			},
		},
		Pricing: []PricingTier{
			{
				Name:  "Starter",
				Price: "IDR 7.500.000",
				WhatYouGet: []string{
					"Landing page 1 halaman",
					"Optimasi konten profil",
					"Integrasi formulir kontak",
				},
			},
			{
				Name:  "Professional",
				Price: "IDR 15.000.000",
				WhatYouGet: []string{
					"Website multi-halaman",
					"Integrasi blog & CMS",
					"Chat assistant dasar",
				},
			},
			{
				Name:  "Elite",
				Price: "Mulai IDR 25.000.000",
				WhatYouGet: []string{
					"AI chat assistant custom",
					"Dashboard analitik",
					"Onboarding & training",
				},
			},
		},
		Contact: ContactInfo{
			Email:    "hello@tany.ai",
			Website:  "https://tany.ai",
			WhatsApp: "+62-811-2345-678",
			LinkedIn: "https://linkedin.com/in/tanyai",
		},
	}
}
