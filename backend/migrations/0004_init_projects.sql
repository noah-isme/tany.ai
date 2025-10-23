-- +migrate Up
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT,
    tech_stack TEXT[] DEFAULT ARRAY[]::TEXT[],
    image_url TEXT,
    project_url TEXT,
    category TEXT,
    "order" INT NOT NULL DEFAULT 0,
    is_featured BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_projects_featured_order ON projects (is_featured DESC, "order");

-- +migrate Down
DROP INDEX IF EXISTS idx_projects_featured_order;
DROP TABLE IF EXISTS projects;
