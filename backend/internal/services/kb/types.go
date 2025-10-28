package kb

import "time"

// Profile represents the public profile returned in the knowledge base payload.
type Profile struct {
	Name      string    `json:"name"`
	Title     string    `json:"title,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Location  string    `json:"location,omitempty"`
	AvatarURL string    `json:"avatarUrl,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// Skill represents a single ordered skill tag.
type Skill struct {
	Name string `json:"name"`
}

// Service encapsulates active service offerings that can be surfaced in chat answers.
type Service struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	DurationLabel string   `json:"durationLabel,omitempty"`
	PriceRange    []string `json:"priceRange,omitempty"`
	Order         int      `json:"order"`
}

// Project highlights past portfolio entries.
type Project struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description,omitempty"`
	TechStack     []string `json:"techStack"`
	ProjectURL    string   `json:"projectUrl,omitempty"`
	Category      string   `json:"category,omitempty"`
	DurationLabel string   `json:"durationLabel,omitempty"`
	PriceLabel    string   `json:"priceLabel,omitempty"`
	BudgetLabel   string   `json:"budgetLabel,omitempty"`
	IsFeatured    bool     `json:"isFeatured"`
	Order         int      `json:"order"`
}

// KnowledgeBase aggregates all public knowledge powering the assistant.
type KnowledgeBase struct {
	Profile  Profile   `json:"profile"`
	Skills   []Skill   `json:"skills"`
	Services []Service `json:"services"`
	Projects []Project `json:"projects"`
}
