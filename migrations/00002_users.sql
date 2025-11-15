-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(255) PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Comments
COMMENT ON TABLE users IS 'System users - team members';
COMMENT ON COLUMN users.id IS 'Unique user identifier (external ID)';
COMMENT ON COLUMN users.username IS 'User display name';
COMMENT ON COLUMN users.team_id IS 'Reference to team the user belongs to';
COMMENT ON COLUMN users.is_active IS 'Active status flag (only active users can be assigned as reviewers)';
COMMENT ON COLUMN users.created_at IS 'Timestamp when user was created';
COMMENT ON COLUMN users.updated_at IS 'Timestamp when user was last updated';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
