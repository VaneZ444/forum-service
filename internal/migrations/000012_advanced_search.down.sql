-- posts: remove full-text search column and index
DROP INDEX IF EXISTS idx_posts_search_vector;
ALTER TABLE posts DROP COLUMN IF EXISTS search_vector;

-- topics: remove full-text search column and index
DROP INDEX IF EXISTS idx_topics_search_vector;
ALTER TABLE topics DROP COLUMN IF EXISTS search_vector;
