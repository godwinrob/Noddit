-- Initial database schema for Noddit

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(32) NOT NULL,
    salt VARCHAR(256) NOT NULL,
    role VARCHAR(255) NOT NULL DEFAULT 'user',
    avatar_address VARCHAR(200),
    first_name VARCHAR(20),
    last_name VARCHAR(20),
    email_address VARCHAR(30),
    join_date DATE
);

-- Subnoddits table
CREATE TABLE IF NOT EXISTS subnoddits (
    sn_id SERIAL PRIMARY KEY,
    sn_name VARCHAR(30) NOT NULL UNIQUE,
    sn_description VARCHAR(200) NOT NULL
);

-- Posts table (supports nested comments)
CREATE TABLE IF NOT EXISTS posts (
    post_id SERIAL PRIMARY KEY,
    parent_post_id BIGINT,
    sn_id BIGINT NOT NULL REFERENCES subnoddits(sn_id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(100) NOT NULL,
    body VARCHAR(2000) NOT NULL,
    image_address VARCHAR,
    post_score BIGINT DEFAULT 1,
    top_level_id BIGINT,
    created_date TIMESTAMP(6) NOT NULL
);

-- Post votes table
CREATE TABLE IF NOT EXISTS post_votes (
    post_id BIGINT NOT NULL REFERENCES posts(post_id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    vote VARCHAR(8) NOT NULL,
    PRIMARY KEY (post_id, user_id)
);

-- Favorites table (polymorphic - favorites both posts and subnoddits)
CREATE TABLE IF NOT EXISTS favorites (
    user_id BIGINT REFERENCES users(id),
    sn_id BIGINT REFERENCES subnoddits(sn_id),
    post_id BIGINT REFERENCES posts(post_id)
);

-- Moderators table
CREATE TABLE IF NOT EXISTS mod (
    sn_id BIGINT NOT NULL REFERENCES subnoddits(sn_id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    PRIMARY KEY (sn_id, user_id)
);

-- Insert default subnoddit (required for user registration)
INSERT INTO subnoddits (sn_name, sn_description)
VALUES ('Cats', 'A home for cats and cat accessories')
ON CONFLICT (sn_name) DO NOTHING;
