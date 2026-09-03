CREATE TABLE IF NOT EXISTS conversation_records (
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    api_key_id BIGINT,
    group_id BIGINT,
    account_id BIGINT,
    provider TEXT NOT NULL DEFAULT '',
    endpoint TEXT NOT NULL DEFAULT '',
    requested_model TEXT NOT NULL DEFAULT '',
    upstream_model TEXT NOT NULL DEFAULT '',
    prompt_content TEXT NOT NULL DEFAULT '',
    response_content TEXT NOT NULL DEFAULT '',
    stream BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT NOT NULL DEFAULT 'unknown',
    complete BOOLEAN NOT NULL DEFAULT FALSE,
    truncated BOOLEAN NOT NULL DEFAULT FALSE,
    input_tokens BIGINT NOT NULL DEFAULT 0,
    output_tokens BIGINT NOT NULL DEFAULT 0,
    duration_ms BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS conversation_records_request_id_uidx ON conversation_records(request_id);
CREATE INDEX IF NOT EXISTS conversation_records_created_at_idx ON conversation_records(created_at DESC);
CREATE INDEX IF NOT EXISTS conversation_records_user_created_idx ON conversation_records(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS conversation_records_api_key_created_idx ON conversation_records(api_key_id, created_at DESC);
CREATE INDEX IF NOT EXISTS conversation_records_model_created_idx ON conversation_records(requested_model, created_at DESC);
