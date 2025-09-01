ALTER TABLE posts
    DROP COLUMN IF EXISTS views_count,
    DROP COLUMN IF EXISTS comments_count,
    DROP COLUMN IF EXISTS likes_count;