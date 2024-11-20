CREATE TABLE IF NOT EXISTS user_invitations (
    token bytea PRIMARY KEY NOT NULL ,
    user_id bigint NOT NULL

)