
-- Enable the uuid-ossp extension if it's not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email text,
    password text,
    subscription_id text,
    subscription_type text,
    channels text[]
);
CREATE TABLE files (
                       id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       bucket_name text,
                       type text,
                       post_id uuid,
                       FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
CREATE TABLE posts (
                      id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                      text text,
                      channel_name text,
                      type text,
                      user_id UUID
);
CREATE TABLE scheduled_posts (
                      id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                      time timestamp,
                      post_id UUID,
                      FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
