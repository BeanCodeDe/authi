# Authi
A lightweight authentication service written in go

![Build](https://img.shields.io/github/workflow/status/BeanCodeDe/authi/MainPipeline.svg)
![License](https://img.shields.io/github/license/BeanCodeDe/authi.svg)
![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)
![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/BeanCodeDe/authi.svg)

## About
Authi is a lightweight authentication service written in go. It covers the basic use cases of creating and deleting users as well as refresh tokens and update passwords.

## Usage

### Configuration

The application can be started with different environment variables to configure specific behavior.


| Name of environment variable | Description                                                           | Mandatory          | Default         |
|:-----------------------------|:----------------------------------------------------------------------|:-------------------|:----------------|
| LOG_LEVEL                    | Log level of console output. You can choose between debug, info, warn | :x:                | info            |
| ADDRESS                      | Server address on that Authi runs                                     | :x:                | 0.0.0.0         |
| PORT                         | Server port on that Authi runs                                        | :x:                | 1203            |
| PRIVATE_KEY_PATH             | Path to the RSA private key file for signing jwt tokens               | :heavy_check_mark: | -               |
| PUBLIC_KEY_PATH              | Path to the RSA public key to validate jwt tokens                     | :heavy_check_mark: | -               |
| DATABASE                     | Used database to store user data                                      | :x:                | postgresql      |
| POSTGRES_USER                | User of postgres database                                             | :x:                | postgres        |
| POSTGRES_PASSWORD            | Password of postgres database                                         | :heavy_check_mark: | -               |
| POSTGRES_DB                  | Database name that should be used in postgres                         | :x:                | postgres        |
| POSTGRES_HOST                | Server address of Postgres database                                   | :x:                | postgres        |
| POSTGRES_PORT                | Server port oft Postgres database                                     | :x:                | 5432            |
| POSTGRES_OPTIONS             | Connection options of Postgres database                               | :x:                | sslmode=disable |

### RSA-Key



### API Interface
The offered api interfaces can be find in this [Swagger UI](https://beancodede.github.io/authi/) or in the folder [/docs](https://github.com/BeanCodeDe/authi/tree/main/docs).

### Adapter

