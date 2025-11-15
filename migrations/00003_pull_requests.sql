-- +goose Up
-- +goose StatementBegin
CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS pull_requests (
  id VARCHAR(255) PRIMARY KEY,
  name VARCHAR(500) NOT NULL,
  author_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
  status pr_status NOT NULL DEFAULT 'OPEN',
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  merged_at TIMESTAMPTZ DEFAULT NULL,

  -- Ensure merged_at consistency with status
  CONSTRAINT chk_pr_merged_at
    CHECK (
      (status = 'MERGED' AND merged_at IS NOT NULL) OR
      (status = 'OPEN' AND merged_at IS NULL)
    )
);

-- Comments
COMMENT ON TABLE pull_requests IS 'Pull requests for code review';
COMMENT ON COLUMN pull_requests.id IS 'Unique pull request identifier';
COMMENT ON COLUMN pull_requests.name IS 'Pull request title/description';
COMMENT ON COLUMN pull_requests.author_id IS 'ID of the PR author';
COMMENT ON COLUMN pull_requests.status IS 'PR status: OPEN (active) or MERGED';
COMMENT ON COLUMN pull_requests.created_at IS 'Timestamp when PR was created';
COMMENT ON COLUMN pull_requests.updated_at IS 'Timestamp when PR was last updated';
COMMENT ON COLUMN pull_requests.merged_at IS 'Timestamp when PR was merged (NULL for OPEN)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pull_requests;
DROP TYPE IF EXISTS pr_status;
-- +goose StatementEnd
