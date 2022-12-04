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

To run the application probably 

- POSTGRES_USER
- POSTGRES_DB 
- POSTGRES_PASSWORD 
- POSTGRES_HOST 
- POSTGRES_PORT 
- LOG_LEVEL 
- PUBLIC_KEY_PATH 
- PRIVATE_KEY_PATH 

### API Interface
The offered api interfaces can be find in this [Swagger UI](https://beancodede.github.io/authi/) or in the folder [/docs](https://github.com/BeanCodeDe/authi/tree/main/docs).

### Adapter
