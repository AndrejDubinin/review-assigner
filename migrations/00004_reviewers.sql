-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reviewers (
  id SERIAL PRIMARY KEY,
  pull_request_id VARCHAR(255) NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
  user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  assigned_at TIMESTAMPTZ DEFAULT NOW(),
  replaced_at TIMESTAMPTZ NULL,
  is_current BOOLEAN DEFAULT true,

  -- Uniqueness: one current reviewer can only be assigned once per PR
  CONSTRAINT unique_current_reviewer
    UNIQUE (pull_request_id, user_id, is_current)
    DEFERRABLE INITIALLY DEFERRED
);

-- Comments
COMMENT ON TABLE reviewers IS 'Junction table linking pull requests to assigned reviewers';
COMMENT ON COLUMN reviewers.id IS 'Auto-incrementing record identifier';
COMMENT ON COLUMN reviewers.pull_request_id IS 'Reference to pull request';
COMMENT ON COLUMN reviewers.user_id IS 'Reference to assigned reviewer';
COMMENT ON COLUMN reviewers.assigned_at IS 'Timestamp when reviewer was assigned';
COMMENT ON COLUMN reviewers.replaced_at IS 'Timestamp when reviewer was replaced (NULL if current)';
COMMENT ON COLUMN reviewers.is_current IS 'true - current reviewer, false - was replaced';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reviewers;
-- +goose StatementEnd
