CREATE OR REPLACE FUNCTION increment_topic_posts_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE topics SET posts_count = posts_count + 1, last_activity = NOW()
        WHERE id = NEW.topic_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE topics SET posts_count = GREATEST(posts_count - 1, 0), last_activity = NOW()
        WHERE id = OLD.topic_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_posts_count
AFTER INSERT OR DELETE ON posts
FOR EACH ROW
EXECUTE FUNCTION increment_topic_posts_count();


CREATE OR REPLACE FUNCTION increment_category_topics_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE categories SET topics_count = topics_count + 1
        WHERE id = NEW.category_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE categories SET topics_count = GREATEST(topics_count - 1, 0)
        WHERE id = OLD.category_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_topics_count
AFTER INSERT OR DELETE ON topics
FOR EACH ROW
EXECUTE FUNCTION increment_category_topics_count();


CREATE OR REPLACE FUNCTION increment_post_comments_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE posts SET comments_count = comments_count + 1
        WHERE id = NEW.post_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE posts SET comments_count = GREATEST(comments_count - 1, 0)
        WHERE id = OLD.post_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_comments_count
AFTER INSERT OR DELETE ON comments
FOR EACH ROW
EXECUTE FUNCTION increment_post_comments_count();
