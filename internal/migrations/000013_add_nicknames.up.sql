ALTER TABLE topics
ADD COLUMN author_nickname VARCHAR(255);

ALTER TABLE posts
ADD COLUMN author_nickname VARCHAR(255);

ALTER TABLE comments
ADD COLUMN author_nickname VARCHAR(255);