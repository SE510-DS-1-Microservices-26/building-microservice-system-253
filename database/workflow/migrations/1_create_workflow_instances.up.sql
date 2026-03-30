CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

CREATE TABLE IF NOT EXISTS workflow_instances
(
    workflow_id UUID         PRIMARY KEY DEFAULT uuid_generate_v7(),
    type        VARCHAR(100) NOT NULL,
    state       VARCHAR(100) NOT NULL,
    payload     JSONB        NOT NULL DEFAULT '{}',
    last_error  TEXT,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_workflow_instances_type  ON workflow_instances (type);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_state ON workflow_instances (state);
