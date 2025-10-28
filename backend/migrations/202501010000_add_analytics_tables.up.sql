CREATE TABLE IF NOT EXISTS analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    event_type TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'unknown',
    provider TEXT NOT NULL DEFAULT 'unknown',
    duration_ms INT NOT NULL DEFAULT 0,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    user_agent TEXT,
    metadata JSONB NOT NULL DEFAULT '{}'::JSONB
);

CREATE INDEX IF NOT EXISTS idx_analytics_events_timestamp ON analytics_events (timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_events_event_type ON analytics_events (event_type);
CREATE INDEX IF NOT EXISTS idx_analytics_events_provider ON analytics_events (provider);
CREATE INDEX IF NOT EXISTS idx_analytics_events_source ON analytics_events (source);

CREATE TABLE IF NOT EXISTS analytics_summary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date DATE NOT NULL UNIQUE,
    total_chats INT NOT NULL DEFAULT 0,
    avg_response_time DOUBLE PRECISION NOT NULL DEFAULT 0,
    success_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
    unique_users INT NOT NULL DEFAULT 0,
    conversions INT NOT NULL DEFAULT 0,
    provider_breakdown JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_analytics_summary_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_analytics_summary_updated_at ON analytics_summary;
CREATE TRIGGER trg_analytics_summary_updated_at
BEFORE UPDATE ON analytics_summary
FOR EACH ROW EXECUTE FUNCTION update_analytics_summary_updated_at();
