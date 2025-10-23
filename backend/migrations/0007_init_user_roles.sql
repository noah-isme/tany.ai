-- +migrate Up
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles (role);

-- +migrate Down
DROP INDEX IF EXISTS idx_user_roles_role;
DROP TABLE IF EXISTS user_roles;
