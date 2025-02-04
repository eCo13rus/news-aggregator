CREATE TABLE IF NOT EXISTS posts (
                                     id SERIAL PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     content TEXT NOT NULL,
                                     pub_time BIGINT NOT NULL,
                                     link TEXT NOT NULL UNIQUE,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_pub_time ON posts(pub_time DESC);
CREATE INDEX IF NOT EXISTS idx_posts_link ON posts(link);
CREATE OR REPLACE FUNCTION update_updated_at_column()

RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();