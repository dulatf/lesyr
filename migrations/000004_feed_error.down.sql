ALTER TABLE feeds DROP COLUMN last_error IF EXISTS;

DROP INDEX idx_feeds_last_fetched_at IF EXISTS;