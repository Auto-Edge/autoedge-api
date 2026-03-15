CREATE TABLE IF NOT EXISTS models (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS model_versions (
    id TEXT PRIMARY KEY,
    model_id TEXT NOT NULL REFERENCES models(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size_bytes BIGINT,
    file_hash TEXT,
    precision TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_published BOOLEAN DEFAULT FALSE,
    download_count INTEGER DEFAULT 0
);

-- Optional: Index for faster lookups
CREATE INDEX idx_models_created_at ON models(created_at DESC);
CREATE INDEX idx_model_versions_model_id ON model_versions(model_id);