ALTER TABLE post_tags
    DROP CONSTRAINT IF EXISTS post_tags_post_id_fkey;

ALTER TABLE post_tags
    ADD CONSTRAINT post_tags_post_id_fkey
        FOREIGN KEY (post_id) REFERENCES posts(id);

ALTER TABLE comments
    DROP CONSTRAINT IF EXISTS comments_post_id_fkey;

ALTER TABLE comments
    ADD CONSTRAINT comments_post_id_fkey
        FOREIGN KEY (post_id) REFERENCES posts(id);