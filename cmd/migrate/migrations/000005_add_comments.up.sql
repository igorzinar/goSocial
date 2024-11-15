CREATE TABLE IF not exists comments (
   id BIGSERIAL PRIMARY KEY ,
    post_id BIGSERIAL NOT NULL,
    user_id BIGSERIAL NOT NULL ,
    content TEXT NOT NULL ,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
)