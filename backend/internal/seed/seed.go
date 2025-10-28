package seed

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tanydotai/tanyai/backend/internal/auth"
)

//go:embed data/*.json
var embeddedData embed.FS

var defaultDataFS fs.FS

func init() {
	sub, err := fs.Sub(embeddedData, "data")
	if err != nil {
		panic(fmt.Errorf("seed: failed to initialize embedded data: %w", err))
	}
	defaultDataFS = sub
}

type profileSeed struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Title     string    `db:"title" json:"title"`
	Bio       string    `db:"bio" json:"bio"`
	Email     string    `db:"email" json:"email"`
	Phone     string    `db:"phone" json:"phone"`
	Location  string    `db:"location" json:"location"`
	AvatarURL string    `db:"avatar_url" json:"avatar_url"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type skillSeed struct {
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Order int    `db:"order" json:"order"`
}

type serviceSeed struct {
	ID            string  `db:"id" json:"id"`
	Name          string  `db:"name" json:"name"`
	Description   string  `db:"description" json:"description"`
	PriceMin      float64 `db:"price_min" json:"price_min"`
	PriceMax      float64 `db:"price_max" json:"price_max"`
	Currency      string  `db:"currency" json:"currency"`
	DurationLabel string  `db:"duration_label" json:"duration_label"`
	IsActive      bool    `db:"is_active" json:"is_active"`
	Order         int     `db:"order" json:"order"`
}

type projectSeed struct {
	ID            string `db:"id" json:"id"`
	Title         string `db:"title" json:"title"`
	Description   string `db:"description" json:"description"`
	TechStack     string `db:"tech_stack" json:"-"`
	ImageURL      string `db:"image_url" json:"image_url"`
	ProjectURL    string `db:"project_url" json:"project_url"`
	Category      string `db:"category" json:"category"`
	DurationLabel string `db:"duration_label" json:"duration_label"`
	PriceLabel    string `db:"price_label" json:"price_label"`
	BudgetLabel   string `db:"budget_label" json:"budget_label"`
	Order         int    `db:"order" json:"order"`
	IsFeatured    bool   `db:"is_featured" json:"is_featured"`
}

type leadSeed struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Message   string    `db:"message" json:"message"`
	Source    string    `db:"source" json:"source"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

type adminUserSeed struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type projectSeedFile struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	TechStack     []string `json:"tech_stack"`
	ImageURL      string   `json:"image_url"`
	ProjectURL    string   `json:"project_url"`
	Category      string   `json:"category"`
	DurationLabel string   `json:"duration_label"`
	PriceLabel    string   `json:"price_label"`
	BudgetLabel   string   `json:"budget_label"`
	Order         int      `json:"order"`
	IsFeatured    bool     `json:"is_featured"`
}

type seedPayload struct {
	Profile  profileSeed
	Skills   []skillSeed
	Services []serviceSeed
	Projects []projectSeed
	Leads    []leadSeed
	Admin    adminUserSeed
}

type Seeder struct {
	db *sqlx.DB
	fs fs.FS
}

func NewSeeder(db *sqlx.DB) *Seeder {
	fsys := defaultDataFS
	if override := os.Getenv("SEED_DATA_PATH"); override != "" {
		fsys = os.DirFS(override)
	}
	return &Seeder{db: db, fs: fsys}
}

func (s *Seeder) Seed(ctx context.Context) error {
	payload, err := loadSeedData(s.fs)
	if err != nil {
		return err
	}
	return seedAll(ctx, s.db, payload)
}

func Seed(ctx context.Context, db *sqlx.DB) error {
	return NewSeeder(db).Seed(ctx)
}

func loadSeedData(fsys fs.FS) (seedPayload, error) {
	var payload seedPayload

	var profileFile profileSeed
	if err := readJSONFile(fsys, "profile.json", &profileFile); err != nil {
		return payload, fmt.Errorf("read profile seed: %w", err)
	}
	payload.Profile = profileFile

	if err := readJSONFile(fsys, "skills.json", &payload.Skills); err != nil {
		return payload, fmt.Errorf("read skills seed: %w", err)
	}

	if err := readJSONFile(fsys, "services.json", &payload.Services); err != nil {
		return payload, fmt.Errorf("read services seed: %w", err)
	}

	var projectFiles []projectSeedFile
	if err := readJSONFile(fsys, "projects.json", &projectFiles); err != nil {
		return payload, fmt.Errorf("read projects seed: %w", err)
	}
	payload.Projects = make([]projectSeed, len(projectFiles))
	for i, p := range projectFiles {
		payload.Projects[i] = projectSeed{
			ID:            p.ID,
			Title:         p.Title,
			Description:   p.Description,
			TechStack:     formatTextArray(p.TechStack),
			ImageURL:      p.ImageURL,
			ProjectURL:    p.ProjectURL,
			Category:      p.Category,
			DurationLabel: p.DurationLabel,
			PriceLabel:    p.PriceLabel,
			BudgetLabel:   p.BudgetLabel,
			Order:         p.Order,
			IsFeatured:    p.IsFeatured,
		}
	}

	if err := readJSONFile(fsys, "leads.json", &payload.Leads); err != nil {
		return payload, fmt.Errorf("read leads seed: %w", err)
	}

	if err := readJSONFile(fsys, "admin_user.json", &payload.Admin); err != nil {
		return payload, fmt.Errorf("read admin user seed: %w", err)
	}

	return payload, nil
}

func readJSONFile(fsys fs.FS, name string, dest interface{}) error {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, dest); err != nil {
		return fmt.Errorf("unmarshal %s: %w", name, err)
	}
	return nil
}

func seedAll(ctx context.Context, db *sqlx.DB, payload seedPayload) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := seedProfile(ctx, tx, payload.Profile); err != nil {
		return err
	}
	if err := seedSkills(ctx, tx, payload.Skills); err != nil {
		return err
	}
	if err := seedServices(ctx, tx, payload.Services); err != nil {
		return err
	}
	if err := seedProjects(ctx, tx, payload.Projects); err != nil {
		return err
	}
	if err := seedLeads(ctx, tx, payload.Leads); err != nil {
		return err
	}
	if err := seedAdminUser(ctx, tx, payload.Admin); err != nil {
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
INSERT INTO projects (id, title, description, tech_stack, image_url, project_url, category, duration_label, price_label, budget_label, "order", is_featured)
VALUES (:id, :title, :description, CAST(:tech_stack AS TEXT[]), :image_url, :project_url, :category, :duration_label, :price_label, :budget_label, :order, :is_featured)
ON CONFLICT (id) DO UPDATE SET
        title = EXCLUDED.title,
        description = EXCLUDED.description,
        tech_stack = EXCLUDED.tech_stack,
        image_url = EXCLUDED.image_url,
        project_url = EXCLUDED.project_url,
        category = EXCLUDED.category,
        duration_label = EXCLUDED.duration_label,
        price_label = EXCLUDED.price_label,
        budget_label = EXCLUDED.budget_label,
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
		item.CreatedAt = time.Now().UTC()
		if _, err := tx.NamedExecContext(ctx, query, item); err != nil {
			return err
		}
	}
	return nil
}

func seedAdminUser(ctx context.Context, tx *sqlx.Tx, admin adminUserSeed) error {
	hashed, err := auth.HashPassword(admin.Password)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	params := map[string]interface{}{
		"id":            admin.ID,
		"email":         admin.Email,
		"password_hash": hashed,
		"name":          admin.Name,
	}

	query := `
INSERT INTO users (id, email, password_hash, name)
VALUES (:id, LOWER(:email), :password_hash, :name)
ON CONFLICT (id) DO UPDATE SET
        email = EXCLUDED.email,
        password_hash = EXCLUDED.password_hash,
        name = EXCLUDED.name,
        updated_at = NOW();
`

	if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_roles (user_id, role) VALUES ($1, $2) ON CONFLICT DO NOTHING`, admin.ID, admin.Role)
	return err
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
