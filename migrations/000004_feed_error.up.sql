-- Add last_error column to feeds table
ALTER TABLE feeds ADD COLUMN last_error TEXT;

-- Add index on last_fetched_at for efficient scheduling
CREATE INDEX idx_feeds_last_fetched_at ON feeds(last_fetched_at);