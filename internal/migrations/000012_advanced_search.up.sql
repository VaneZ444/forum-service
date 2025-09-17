-- posts: add full-text search column
ALTER TABLE posts
ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(content, '')), 'B')
) STORED;

CREATE INDEX idx_posts_search_vector ON posts USING GIN(search_vector);

-- topics: add full-text search column
ALTER TABLE topics
ADD COLUMN search_vector tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A')
) STORED;

CREATE INDEX idx_topics_search_vector ON topics USING GIN(search_vector);
