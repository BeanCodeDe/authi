CREATE TABLE auth.user (
    id uuid PRIMARY KEY NOT NULL,
    password varchar NOT NULL,
    salt varchar(32) NOT NULL,
    refresh_token varchar(32),
    refresh_token_expire timestamp,
    created_on timestamp NOT NULL,
    last_login timestamp NOT NULL
);

CREATE INDEX idx_password ON auth.user (password);
CREATE INDEX idx_refresh_token ON auth.user (id, refresh_token, refresh_token_expire);
