CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    "order" INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_skills_order ON skills ("order");

