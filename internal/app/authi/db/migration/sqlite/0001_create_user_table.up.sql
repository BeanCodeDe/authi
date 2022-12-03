CREATE TABLE auth.user (
    id uuid PRIMARY KEY NOT NULL,
    password varchar NOT NULL,
    salt varchar(32) NOT NULL,
    refresh_token varchar(32),
    refresh_token_expire timestamp,
    created_on timestamp NOT NULL,
    last_login timestamp NOT NULL
);
