ALTER TABLE videos ADD CONSTRAINT videos_length_check CHECK (length >= 0);
ALTER TABLE videos ADD CONSTRAINT videos_published_at_check CHECK (published_at <= now());
