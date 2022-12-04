# Authi
A lightweight authentication service written in go

![Build](https://img.shields.io/github/workflow/status/BeanCodeDe/authi/MainPipeline.svg)
![License](https://img.shields.io/github/license/BeanCodeDe/authi.svg)
![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)
![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/BeanCodeDe/authi.svg)

## About
Authi is a lightweight authentication service written in go. It covers the basic use cases of creating and deleting users as well as refresh tokens and update passwords.

---

## Getting started

For a fast setup you can use the following guide:

### Create a folder structure that looks like this:
```
.
├── docker-compose.yml
└── myTokenFolder
```

### Paste into the `docker-compose.yml` the following content

```
version: '3.7'
services:
  postgres:
    image: postgres:latest
    container_name: postgres
    restart: always
    environment: 
      - POSTGRES_PASSWORD=myDatabasePassword
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5 
  authi:
    image: "beancodede/authi:latest"
    container_name: authi
    restart: always
    environment: 
      - POSTGRES_PASSWORD=myDatabasePassword
    ports:
      - 1203:1203
    volumes: 
      - ./myTokenFolder:/token
    depends_on:
      postgres:
        condition: service_healthy
    links:
      - postgres:postgres
```

### Execute the following two commands to generate key's
```
ssh-keygen -t rsa -b 4096 -m PEM -f ./myTokenFolder/jwtRS256.key
```
```
openssl rsa -in ./myTokenFolder/jwtRS256.key -pubout -outform PEM -out ./myTokenFolder/jwtRS256.key.pub
```

>You need the cli programmes `ssh-keygen` and `openssl` in order to execute these commands

### Start docker compose file

```
docker compose up
```

>You need `docker` installed in order to execute this command

---

## Configuration

The application can be started with different environment variables to configure specific behavior.


| Name of environment variable | Description                                                           | Mandatory          | Default         |
|:-----------------------------|:----------------------------------------------------------------------|:-------------------|:----------------|
| LOG_LEVEL                    | Log level of console output. You can choose between debug, info, warn | :x:                | info            |
| ADDRESS                      | Server address on that Authi runs                                     | :x:                | 0.0.0.0         |
| PORT                         | Server port on that Authi runs                                        | :x:                | 1203            |
| PRIVATE_KEY_PATH             | Path to the RSA private key file for signing jwt tokens               | :x:  | -               |
| PUBLIC_KEY_PATH              | Path to the RSA public key to validate jwt tokens                     | :x:  | -               |
| DATABASE                     | Used database to store user data                                      | :x:                | postgresql      |
| POSTGRES_USER                | User of postgres database                                             | :x:                | postgres        |
| POSTGRES_PASSWORD            | Password of postgres database                                         | :heavy_check_mark: | -               |
| POSTGRES_DB                  | Database name that should be used in postgres                         | :x:                | postgres        |
| POSTGRES_HOST                | Server address of Postgres database                                   | :x:                | postgres        |
| POSTGRES_PORT                | Server port oft Postgres database                                     | :x:                | 5432            |
| POSTGRES_OPTIONS             | Connection options of Postgres database                               | :x:                | sslmode=disable |

---

## Adapter

**TODO**

## API Interface
The offered api interfaces can be find in this [Swagger UI](https://beancodede.github.io/authi/) or in the folder [/docs](https://github.com/BeanCodeDe/authi/tree/main/docs).


---

## Setting up project
**TODO**