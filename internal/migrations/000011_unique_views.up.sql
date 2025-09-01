CREATE UNIQUE INDEX IF NOT EXISTS idx_post_views_unique
    ON post_views (post_id, user_id);