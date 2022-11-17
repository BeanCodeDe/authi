CREATE DATABASE authi;
\c authi
CREATE USER authi WITH ENCRYPTED PASSWORD 'secret_password';

CREATE SCHEMA auth; 

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

GRANT CONNECT ON DATABASE authi TO authi;  -- since we revoked from public
GRANT USAGE ON SCHEMA auth TO authi;
GRANT ALL ON ALL TABLES IN SCHEMA auth TO authi;
GRANT ALL ON ALL SEQUENCES IN SCHEMA auth TO authi; -- don't forget those
ALTER DEFAULT PRIVILEGES FOR ROLE authi IN SCHEMA auth
GRANT ALL ON TABLES TO authi;
ALTER DEFAULT PRIVILEGES FOR ROLE authi IN SCHEMA auth
GRANT ALL ON SEQUENCES TO authi;
