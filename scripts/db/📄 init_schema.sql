-- Enable extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Monitors table with auto-generated UUID and created_at
CREATE TABLE IF NOT EXISTS monitors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    interval_seconds INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TYPE monitor_status IF NOT EXISTS AS ENUM ('Up', 'Down');

-- Monitor results table with auto-generated serial id and created_at
CREATE TABLE IF NOT EXISTS monitor_results (
    id SERIAL PRIMARY KEY,
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status monitor_status NOT NULL,
    latency_ms INTEGER,
    response_code INTEGER,
    error TEXT,
    region TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Helpful indexes
CREATE INDEX IF NOT EXISTS idx_monitor_results_monitor_id ON monitor_results(monitor_id);
CREATE INDEX IF NOT EXISTS idx_monitor_results_checked_at ON monitor_results(checked_at);
