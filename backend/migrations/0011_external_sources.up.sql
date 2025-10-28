CREATE TABLE IF NOT EXISTS external_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    base_url TEXT NOT NULL,
    source_type TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    etag TEXT,
    last_modified TIMESTAMPTZ,
    last_synced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (name <> ''),
    CHECK (base_url <> ''),
    CHECK (source_type <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_external_sources_name ON external_sources (LOWER(name));
CREATE UNIQUE INDEX IF NOT EXISTS idx_external_sources_base_url ON external_sources (LOWER(base_url));
CREATE INDEX IF NOT EXISTS idx_external_sources_enabled ON external_sources (enabled);

CREATE TABLE IF NOT EXISTS external_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES external_sources(id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    summary TEXT,
    content TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB,
    published_at TIMESTAMPTZ,
    hash TEXT NOT NULL,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (kind <> ''),
    CHECK (title <> ''),
    CHECK (url <> ''),
    CHECK (hash <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_external_items_source_url ON external_items (source_id, LOWER(url));
CREATE UNIQUE INDEX IF NOT EXISTS idx_external_items_source_hash ON external_items (source_id, hash);
CREATE INDEX IF NOT EXISTS idx_external_items_published_at ON external_items (published_at DESC NULLS LAST);
CREATE INDEX IF NOT EXISTS idx_external_items_kind_visible ON external_items (kind, visible);
