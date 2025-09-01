DROP TRIGGER IF EXISTS trg_posts_count ON posts;
DROP FUNCTION IF EXISTS increment_topic_posts_count;

DROP TRIGGER IF EXISTS trg_topics_count ON topics;
DROP FUNCTION IF EXISTS increment_category_topics_count;

DROP TRIGGER IF EXISTS trg_comments_count ON comments;
DROP FUNCTION IF EXISTS increment_post_comments_count;
