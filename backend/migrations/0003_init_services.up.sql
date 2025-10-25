CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    price_min NUMERIC(12,2),
    price_max NUMERIC(12,2),
    currency TEXT,
    duration_label TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    "order" INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_services_active_order ON services (is_active, "order");

