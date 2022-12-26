# Authi
A lightweight authentication service written in go

![Build](https://img.shields.io/github/workflow/status/BeanCodeDe/authi/Audit_And_Deploy.svg)
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

>**Note**
>You need the cli programmes `ssh-keygen` and `openssl` in order to execute these commands

### Start docker compose file

```
docker compose up
```

>**Note**
>You need `docker` installed in order to execute this command

---

## Configuration

The application can be started with different environment variables to configure specific behavior.


| Name of environment variable | Description                                                           | Mandatory          | Default                 |
|:-----------------------------|:----------------------------------------------------------------------|:-------------------|:------------------------|
| LOG_LEVEL                    | Log level of console output. You can choose between debug, info, warn | :x:                | info                    |
| ADDRESS                      | Server address on that Authi runs                                     | :x:                | 0.0.0.0                 |
| PORT                         | Server port on that Authi runs                                        | :x:                | 1203                    |
| PRIVATE_KEY_PATH             | Path to the RSA private key file for signing jwt tokens               | :x:                | /token/jwtRS256.key     |
| PUBLIC_KEY_PATH              | Path to the RSA public key to validate jwt tokens                     | :x:                | /token/jwtRS256.key.pub |
| DATABASE                     | Used database to store user data                                      | :x:                | postgresql              |
| POSTGRES_USER                | User of postgres database                                             | :x:                | postgres                |
| POSTGRES_PASSWORD            | Password of postgres database                                         | :heavy_check_mark: | -                       |
| POSTGRES_DB                  | Database name that should be used in postgres                         | :x:                | postgres                |
| POSTGRES_HOST                | Server address of Postgres database                                   | :x:                | postgres                |
| POSTGRES_PORT                | Server port oft Postgres database                                     | :x:                | 5432                    |
| POSTGRES_OPTIONS             | Connection options of Postgres database                               | :x:                | sslmode=disable         |
| ACCESS_TOKEN_EXPIRE_TIME     | Time in minutes till an access token is no longer valid               | :x:                | 5                       |
| REFRESH_TOKEN_EXPIRE_TIME    | Time in minutes till an refresh token is no longer valid              | :x:                | 10                      |

---

## API Interface
The offered api interfaces can be find in this [Swagger UI](https://beancodede.github.io/authi/) or in the folder [/docs](https://github.com/BeanCodeDe/authi/tree/main/docs).

## Adapter

To access the authi service from other go echo application, you can use the methods within the adapter package. Therefore two methods are provided:

```Go
GetToken(userId string, password string) (*TokenResponseDTO, error)
```
to get an access and refresh token with your userId and password.


```Go
RefreshToken(userId string, token string, refreshToken string) (*TokenResponseDTO, error)
```
to refresh your token for a logged in user without passing the password again.

>**Note**
>You can only refresh tokens as long as your token and refresh token are valid. 


To initial an authi adapter you have to use the method `NewAuthiAdapter()` within the adapter package.

The whole code could look like this:

```Go
func AdapterExample() {

	//User id that were previously created over REST
	userId := "693227c8-4178-4e72-b3b7-a8b8bae36f1b"

	//Initialize authi adapter
	authiAdapter := adapter.NewAuthiAdapter()

	//Logging in previously created user with password `mySecretUserPassword`
	token, err := authiAdapter.GetToken(userId, "mySecretUserPassword")

	//Checking if an error occurred while loading user token
	if err != nil {
		panic(err)
	}

	//Printing access token for example
	fmt.Println(token.AccessToken)

	//Printing refresh token for example
	fmt.Println(token.RefreshToken)

	//Refreshing tokens to avoid outdated tokens
	refreshedToken, err := authiAdapter.RefreshToken(userId, token.AccessToken, token.RefreshToken)

	//Checking if an error occurred while loading refreshed token
	if err != nil {
		panic(err)
	}

	//Printing refreshed access token for example
	fmt.Println(refreshedToken.AccessToken)

	//Printing refreshed refresh token for example
	fmt.Println(refreshedToken.RefreshToken)
}
```

## Middleware

If you write an echo go application, you also have the possibility to use the already attached middleware. Therefore you can orientate on the following example code:


```Go
func MiddlewareExample() {
	//Initialize parser to validate Tokens
	tokenParser, err := parser.NewJWTParser()

	//Checking if an error occurred while loading jwt parser
	if err != nil {
		panic(err)
	}

	//Initialize middleware
	echoMiddleware := middleware.NewEchoMiddleware(tokenParser)

	//Initialize echo
	e := echo.New()

	//Secure endpoint with method `echoMiddleware.CheckToken`
	e.GET(
		"/someEndpoint",
		func(c echo.Context) error { return c.NoContent(201) },
		echoMiddleware.CheckToken,
	)
}
```

>**Note**
>In order for the code to work properly, the environment variable `PUBLIC_KEY_PATH` must point to the appropriate public key of the Authi server


While using the middleware the following errors could occur:

| Error                    | Description                                   | Instruction                                                                                                                                   |
|:-------------------------|:----------------------------------------------|:----------------------------------------------------------------------------------------------------------------------------------------------|
| ErrTokenNotFound         | Token was not found in request                | Does the passed token starts with "Bearer "; Is the token set in header `Authorization` inside the request                                    |
| ErrClaimCouldNotBeParsed | Token is valid but Claim format is unexpected | Are you using compatible version? Otherwise please create an issue                                                                            |
| ErrTokenNotValid         | Token is not valid                            | A different key file may be used here for checking than is used in the Authi service. Otherwise, the token may have expired.                  |
| ErrWhileReadingKey       | Key file couldn't be read                     | Maybe the file is not on the correct position. Check your environment variable `PUBLIC_KEY_PATH`                                              |
| ErrWhileParsingKey       | Key file couldn't be parsed                   | Maybe the file has not the correct format. Use the commands from `Execute the following two commands to generate key's` to generate key files |

---