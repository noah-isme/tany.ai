package seed

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type profileSeed struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Title     string    `db:"title"`
	Bio       string    `db:"bio"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	Location  string    `db:"location"`
	AvatarURL string    `db:"avatar_url"`
	UpdatedAt time.Time `db:"updated_at"`
}

type skillSeed struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Order int    `db:"order"`
}

type serviceSeed struct {
	ID            string  `db:"id"`
	Name          string  `db:"name"`
	Description   string  `db:"description"`
	PriceMin      float64 `db:"price_min"`
	PriceMax      float64 `db:"price_max"`
	Currency      string  `db:"currency"`
	DurationLabel string  `db:"duration_label"`
	IsActive      bool    `db:"is_active"`
	Order         int     `db:"order"`
}

type projectSeed struct {
	ID          string `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	TechStack   string `db:"tech_stack"`
	ImageURL    string `db:"image_url"`
	ProjectURL  string `db:"project_url"`
	Category    string `db:"category"`
	Order       int    `db:"order"`
	IsFeatured  bool   `db:"is_featured"`
}

type leadSeed struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Message   string    `db:"message"`
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
}

var (
	profileData = profileSeed{
		ID:        "11111111-1111-1111-1111-111111111111",
		Name:      "John Doe",
		Title:     "Full Stack Developer & AI Consultant",
		Bio:       "Freelance developer yang berfokus pada pembuatan aplikasi web modern, integrasi AI, dan automasi bisnis.",
		Email:     "john@example.com",
		Phone:     "+62812345678",
		Location:  "Jakarta, Indonesia",
		AvatarURL: "https://images.example.com/john-doe-avatar.png",
		UpdatedAt: time.Now().UTC(),
	}

	skillData = []skillSeed{
		{ID: "22222222-1111-1111-1111-111111111111", Name: "React.js", Order: 1},
		{ID: "22222222-2222-1111-1111-111111111111", Name: "Golang", Order: 2},
		{ID: "22222222-3333-1111-1111-111111111111", Name: "PostgreSQL", Order: 3},
		{ID: "22222222-4444-1111-1111-111111111111", Name: "UI/UX Design", Order: 4},
		{ID: "22222222-5555-1111-1111-111111111111", Name: "Prompt Engineering", Order: 5},
		{ID: "22222222-6666-1111-1111-111111111111", Name: "Automation & No-Code", Order: 6},
	}

	serviceData = []serviceSeed{
		{
			ID:            "33333333-1111-1111-1111-111111111111",
			Name:          "Website Development",
			Description:   "Pembuatan website custom dengan stack modern dan optimasi performa.",
			PriceMin:      5000000,
			PriceMax:      20000000,
			Currency:      "IDR",
			DurationLabel: "3-6 minggu",
			IsActive:      true,
			Order:         1,
		},
		{
			ID:            "33333333-2222-1111-1111-111111111111",
			Name:          "Landing Page Sprint",
			Description:   "Landing page responsif untuk kampanye marketing atau produk baru.",
			PriceMin:      2500000,
			PriceMax:      7500000,
			Currency:      "IDR",
			DurationLabel: "1-2 minggu",
			IsActive:      true,
			Order:         2,
		},
		{
			ID:            "33333333-3333-1111-1111-111111111111",
			Name:          "AI Chatbot Integration",
			Description:   "Integrasi chatbot AI berbasis GPT yang dikustomisasi dengan knowledge base bisnis Anda.",
			PriceMin:      6000000,
			PriceMax:      18000000,
			Currency:      "IDR",
			DurationLabel: "2-3 minggu",
			IsActive:      true,
			Order:         3,
		},
		{
			ID:            "33333333-4444-1111-1111-111111111111",
			Name:          "Brand & UI Refresh",
			Description:   "Re-desain identitas visual dan UI dashboard agar lebih modern dan konsisten.",
			PriceMin:      3500000,
			PriceMax:      12000000,
			Currency:      "IDR",
			DurationLabel: "2-4 minggu",
			IsActive:      true,
			Order:         4,
		},
		{
			ID:            "33333333-5555-1111-1111-111111111111",
			Name:          "Product Maintenance Retainer",
			Description:   "Paket maintenance bulanan untuk feature updates, bug fix, dan monitoring performa.",
			PriceMin:      1500000,
			PriceMax:      5000000,
			Currency:      "IDR",
			DurationLabel: "Bulanan",
			IsActive:      true,
			Order:         5,
		},
	}

	projectData = []projectSeed{
		{
			ID:          "44444444-1111-1111-1111-111111111111",
			Title:       "E-commerce Fashion Platform",
			Description: "Platform toko online dengan katalog dinamis, pembayaran terintegrasi, dan dashboard analitik.",
			TechStack:   formatTextArray([]string{"Next.js", "Golang", "Supabase", "Midtrans"}),
			ImageURL:    "https://images.example.com/projects/fashion-store.png",
			ProjectURL:  "https://fashion-store.example.com",
			Category:    "E-commerce",
			Order:       1,
			IsFeatured:  true,
		},
		{
			ID:          "44444444-2222-1111-1111-111111111111",
			Title:       "SaaS Analytics Dashboard",
			Description: "Dashboard real-time untuk monitoring KPI bisnis dengan integrasi data warehouse.",
			TechStack:   formatTextArray([]string{"React", "Gin", "PostgreSQL", "Grafana"}),
			ImageURL:    "https://images.example.com/projects/saas-analytics.png",
			ProjectURL:  "https://analytics.example.com",
			Category:    "SaaS",
			Order:       2,
			IsFeatured:  true,
		},
		{
			ID:          "44444444-3333-1111-1111-111111111111",
			Title:       "AI Customer Support Bot",
			Description: "Chatbot AI untuk support pelanggan dengan prompt builder dan retrieval knowledge base.",
			TechStack:   formatTextArray([]string{"Next.js", "Golang", "OpenAI", "Supabase"}),
			ImageURL:    "https://images.example.com/projects/ai-support.png",
			ProjectURL:  "https://ai-support.example.com",
			Category:    "AI & Automation",
			Order:       3,
			IsFeatured:  true,
		},
		{
			ID:          "44444444-4444-1111-1111-111111111111",
			Title:       "Company Profile Revamp",
			Description: "Redesign website company profile dengan storytelling dan ilustrasi custom.",
			TechStack:   formatTextArray([]string{"Astro", "Tailwind", "Framer Motion"}),
			ImageURL:    "https://images.example.com/projects/company-profile.png",
			ProjectURL:  "https://company-profile.example.com",
			Category:    "Creative",
			Order:       4,
			IsFeatured:  false,
		},
	}

	leadData = []leadSeed{
		{
			ID:        "55555555-1111-1111-1111-111111111111",
			Name:      "Jane Startup",
			Email:     "jane@startup.io",
			Message:   "Hi John, kami butuh bantuan membuat MVP SaaS analytics. Bisa kirim proposal?",
			Source:    "landing_page_form",
			CreatedAt: time.Now().UTC(),
		},
	}
)

// Seed performs idempotent upsert operations to populate the core tables with
// representative knowledge base data.
func Seed(ctx context.Context, db *sqlx.DB) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := seedProfile(ctx, tx, profileData); err != nil {
		return err
	}
	if err := seedSkills(ctx, tx, skillData); err != nil {
		return err
	}
	if err := seedServices(ctx, tx, serviceData); err != nil {
		return err
	}
	if err := seedProjects(ctx, tx, projectData); err != nil {
		return err
	}
	if err := seedLeads(ctx, tx, leadData); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func seedProfile(ctx context.Context, tx *sqlx.Tx, data profileSeed) error {
	data.UpdatedAt = time.Now().UTC()
	query := `
INSERT INTO profile (id, name, title, bio, email, phone, location, avatar_url, updated_at)
VALUES (:id, :name, :title, :bio, :email, :phone, :location, :avatar_url, :updated_at)
ON CONFLICT (id) DO UPDATE SET
        name = EXCLUDED.name,
        title = EXCLUDED.title,
        bio = EXCLUDED.bio,
        email = EXCLUDED.email,
        phone = EXCLUDED.phone,
        location = EXCLUDED.location,
        avatar_url = EXCLUDED.avatar_url,
        updated_at = EXCLUDED.updated_at;
`
	_, err := tx.NamedExecContext(ctx, query, data)
	return err
}

func seedSkills(ctx context.Context, tx *sqlx.Tx, data []skillSeed) error {
	query := `
INSERT INTO skills (id, name, "order")
VALUES (:id, :name, :order)
ON CONFLICT (id) DO UPDATE SET
        name = EXCLUDED.name,
        "order" = EXCLUDED."order";
`
	for _, item := range data {
		if _, err := tx.NamedExecContext(ctx, query, item); err != nil {
			return err
		}
	}
	return nil
}

func seedServices(ctx context.Context, tx *sqlx.Tx, data []serviceSeed) error {
	query := `
INSERT INTO services (id, name, description, price_min, price_max, currency, duration_label, is_active, "order")
VALUES (:id, :name, :description, :price_min, :price_max, :currency, :duration_label, :is_active, :order)
ON CONFLICT (id) DO UPDATE SET
        name = EXCLUDED.name,
        description = EXCLUDED.description,
        price_min = EXCLUDED.price_min,
        price_max = EXCLUDED.price_max,
        currency = EXCLUDED.currency,
        duration_label = EXCLUDED.duration_label,
        is_active = EXCLUDED.is_active,
        "order" = EXCLUDED."order";
`
	for _, item := range data {
		if _, err := tx.NamedExecContext(ctx, query, item); err != nil {
			return err
		}
	}
	return nil
}

func seedProjects(ctx context.Context, tx *sqlx.Tx, data []projectSeed) error {
	query := `
INSERT INTO projects (id, title, description, tech_stack, image_url, project_url, category, "order", is_featured)
VALUES (:id, :title, :description, CAST(:tech_stack AS TEXT[]), :image_url, :project_url, :category, :order, :is_featured)
ON CONFLICT (id) DO UPDATE SET
        title = EXCLUDED.title,
        description = EXCLUDED.description,
        tech_stack = EXCLUDED.tech_stack,
        image_url = EXCLUDED.image_url,
        project_url = EXCLUDED.project_url,
        category = EXCLUDED.category,
        "order" = EXCLUDED."order",
        is_featured = EXCLUDED.is_featured;
`
	for _, item := range data {
		if _, err := tx.NamedExecContext(ctx, query, item); err != nil {
			return err
		}
	}
	return nil
}

func seedLeads(ctx context.Context, tx *sqlx.Tx, data []leadSeed) error {
	query := `
INSERT INTO leads (id, name, email, message, source, created_at)
VALUES (:id, :name, :email, :message, :source, :created_at)
ON CONFLICT (id) DO UPDATE SET
        name = EXCLUDED.name,
        email = EXCLUDED.email,
        message = EXCLUDED.message,
        source = EXCLUDED.source,
        created_at = EXCLUDED.created_at;
`
	for _, item := range data {
		if _, err := tx.NamedExecContext(ctx, query, item); err != nil {
			return err
		}
	}
	return nil
}

func formatTextArray(values []string) string {
	if len(values) == 0 {
		return "{}"
	}
	escaped := make([]string, len(values))
	for i, v := range values {
		escaped[i] = strconv.Quote(v)
	}
	return fmt.Sprintf("{%s}", strings.Join(escaped, ","))
}
