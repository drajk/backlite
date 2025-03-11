CREATE TABLE IF NOT EXISTS backlite_tasks (
    id VARCHAR(255) PRIMARY KEY,
    created_at BIGINT NOT NULL,
    queue VARCHAR(255) NOT NULL,
    task LONGBLOB NOT NULL,
    wait_until BIGINT,
    claimed_at BIGINT,
    last_executed_at BIGINT,
    attempts INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS backlite_tasks_completed (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    created_at BIGINT NOT NULL,
    queue VARCHAR(255) NOT NULL,
    last_executed_at BIGINT,
    attempts INT NOT NULL,
    last_duration_micro BIGINT,
    succeeded INT,
    task LONGBLOB,
    expires_at BIGINT,
    error TEXT
);

CREATE INDEX backlite_tasks_wait_until ON backlite_tasks (wait_until);
