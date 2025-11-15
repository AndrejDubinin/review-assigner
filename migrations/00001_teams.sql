-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Comments
COMMENT ON TABLE teams IS 'Development teams';
COMMENT ON COLUMN teams.id IS 'Auto-incrementing team identifier';
COMMENT ON COLUMN teams.name IS 'Unique team name';
COMMENT ON COLUMN teams.created_at IS 'Timestamp when team was created';
COMMENT ON COLUMN teams.updated_at IS 'Timestamp when team was last updated';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
