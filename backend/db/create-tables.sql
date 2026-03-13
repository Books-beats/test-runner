
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE tests (
    id BIGSERIAL PRIMARY KEY,
    
    name VARCHAR(255) NOT NULL,
    
    url TEXT NOT NULL,
    
    method VARCHAR(10) NOT NULL CHECK (method IN ('GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS')),
    
    headers JSONB DEFAULT '{}'::jsonb,
    
    body TEXT,
    
    user_id BIGINT NOT NULL,
    
    expected_response TEXT,

    status_code INT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE test_runs (
    id BIGSERIAL PRIMARY KEY,
    
    test_id BIGINT NOT NULL,
    
    concurrency INT NOT NULL CHECK (concurrency > 0),
    
    total INT,
    
    passed INT,
    
    failed INT,
    
    avg_duration_ms INT,
    min_duration_ms INT,
    max_duration_ms INT,
    
    status VARCHAR(50) DEFAULT 'pending',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT fk_test
        FOREIGN KEY (test_id)
        REFERENCES tests(id)
        ON DELETE CASCADE
);

CREATE TABLE job_results (
    id BIGSERIAL PRIMARY KEY,
    
    test_run_id BIGINT NOT NULL,
    
    job_number INT NOT NULL,
    
    status_code INT,
    
    duration_ms INT,
    
    response_size INT,
    
    passed BOOLEAN NOT NULL DEFAULT FALSE,
    
    error TEXT,

    status TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    completed_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_test_run
        FOREIGN KEY (test_run_id)
        REFERENCES test_runs(id)
        ON DELETE CASCADE
);