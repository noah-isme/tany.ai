-- Enable pgvector extension and create embeddings storage.
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS embeddings (
    id UUID PRIMARY KEY,
    kind TEXT NOT NULL,
    ref_id UUID NULL,
    content TEXT NOT NULL,
    vector VECTOR(1536),
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS embeddings_kind_ref_id_uq
    ON embeddings(kind, ref_id)
    WHERE ref_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS embeddings_kind_idx
    ON embeddings(kind);

CREATE INDEX IF NOT EXISTS embeddings_updated_at_idx
    ON embeddings(updated_at DESC);

CREATE INDEX IF NOT EXISTS embeddings_vector_idx
    ON embeddings
    USING ivfflat (vector)
    WITH (lists = 100);

CREATE TABLE IF NOT EXISTS embedding_config (
    key TEXT PRIMARY KEY,
    value JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO embedding_config (key, value)
VALUES (
    'personalization',
    jsonb_build_object(
        'weight', 0.65,
        'lastReindexedAt', NULL,
        'lastResetAt', NULL
    )
)
ON CONFLICT (key) DO NOTHING;
