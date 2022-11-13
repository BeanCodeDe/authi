--SPACELIGHT-AUTH
CREATE DATABASE authi;
\c authi
CREATE USER authi WITH ENCRYPTED PASSWORD 'secret_password';

CREATE SCHEMA authi; 

CREATE TABLE authi.user (
    id uuid PRIMARY KEY NOT NULL,
    password varchar NOT NULL,
    salt varchar(32) NOT NULL,
    refresh_token varchar,
    created_on timestamp NOT NULL,
    last_login timestamp NOT NULL
);

CREATE INDEX idx_password ON spacelight.user (password);
CREATE INDEX refresh_token ON spacelight.user (password);

GRANT CONNECT ON DATABASE authi TO authi;  
GRANT USAGE ON SCHEMA authi TO authi;
GRANT ALL ON ALL TABLES IN SCHEMA authi TO authi;
GRANT ALL ON ALL SEQUENCES IN SCHEMA authi TO authi;  
ALTER DEFAULT PRIVILEGES FOR ROLE authi IN SCHEMA authi
GRANT ALL ON TABLES TO authi;
ALTER DEFAULT PRIVILEGES FOR ROLE authi IN SCHEMA authi
GRANT ALL ON SEQUENCES TO authi;