CREATE TABLE IF NOT EXISTS backlite_tasks (
    id VARCHAR(255) PRIMARY KEY,
    created_at INT NOT NULL,
    queue VARCHAR(255) NOT NULL,
    task LONGBLOB NOT NULL,
    wait_until INT,
    claimed_at INT,
    last_executed_at INT,
    attempts INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS backlite_tasks_completed (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    created_at INT NOT NULL,
    queue VARCHAR(255) NOT NULL,
    last_executed_at INT,
    attempts INT NOT NULL,
    last_duration_micro INT,
    succeeded INT,
    task LONGBLOB,
    expires_at INT,
    error TEXT
);

CREATE INDEX backlite_tasks_wait_until ON backlite_tasks (wait_until);
