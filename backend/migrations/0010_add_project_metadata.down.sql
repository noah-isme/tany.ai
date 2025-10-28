ALTER TABLE projects
    DROP COLUMN IF EXISTS duration_label,
    DROP COLUMN IF EXISTS price_label,
    DROP COLUMN IF EXISTS budget_label;
