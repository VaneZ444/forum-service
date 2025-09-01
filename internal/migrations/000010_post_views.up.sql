CREATE TABLE post_views (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id BIGINT, -- NULL если анонимный
    viewed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION increment_post_views_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE posts
    SET views_count = views_count + 1
    WHERE id = NEW.post_id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_post_views_count
AFTER INSERT ON post_views
FOR EACH ROW
EXECUTE FUNCTION increment_post_views_count();