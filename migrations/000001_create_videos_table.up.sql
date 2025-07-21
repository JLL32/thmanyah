CREATE TYPE video_type AS ENUM ('podcast', 'documentary');

CREATE TABLE IF NOT EXISTS videos (
   video_id VARCHAR(11) PRIMARY KEY,
   title text NOT NULL,
   description text NOT NULL,
   type video_type NOT NULL,
   length integer NOT NULL,
   language VARCHAR(2) NOT NULL CHECK (language ~ '^[a-z]{2}$'),
   published_at timestamp(0) with time zone NOT NULL,
   created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
   version integer NOT NULL DEFAULT 1
);
