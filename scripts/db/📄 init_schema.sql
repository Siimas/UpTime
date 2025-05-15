-- Enable extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Monitors table with auto-generated UUID and created_at
CREATE TABLE monitors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url TEXT NOT NULL,
    interval_seconds INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Monitor results table with auto-generated serial id and created_at
CREATE TABLE monitor_results (
    id SERIAL PRIMARY KEY,
    monitor_id UUID NOT NULL REFERENCES monitors(id),
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status TEXT NOT NULL,
    latency_ms INTEGER,
    response_code INTEGER,
    error TEXT,
    region TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Helpful indexes
CREATE INDEX idx_monitor_results_monitor_id ON monitor_results(monitor_id);
CREATE INDEX idx_monitor_results_checked_at ON monitor_results(checked_at);
