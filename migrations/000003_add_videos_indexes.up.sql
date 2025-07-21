CREATE INDEX IF NOT EXISTS videos_title_idx ON videos USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS videos_description_idx ON videos USING GIN (to_tsvector('simple', description));
