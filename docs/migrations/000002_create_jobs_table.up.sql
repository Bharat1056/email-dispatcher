CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    queue_name VARCHAR(50) NOT NULL DEFAULT 'default',
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CONSTRAINT chk_job_status CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'dlq')),
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    run_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_attempts (
    id BIGSERIAL PRIMARY KEY,
    job_id BIGINT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    attempt_number INT NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at TIMESTAMP WITH TIME ZONE,
    error_log TEXT,
    duration_ms INT,
    queue_duration_ms INT
);

CREATE INDEX IF NOT EXISTS idx_job_attempts_job_id ON job_attempts(job_id);
