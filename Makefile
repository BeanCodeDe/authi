SRC_PATH?=./cmd/authi
APP_NAME?=authi
DOCKER_COMPOSE_PATH?=./deployments/docker-compose-postgres.yml
DOCKER_PATH?=./build/Dockerfile
ENV_CONFIG?=./deployments/test.env

app.build:
	go mod download
	go build -o $(APP_NAME) $(SRC_PATH)

app.ut.run:
	go test ./internal/... -v

app.jt.run:
	docker compose --env-file $(ENV_CONFIG) --file $(DOCKER_COMPOSE_PATH) up --build --force-recreate -d
	go test ./test
	docker compose --env-file $(ENV_CONFIG) --file $(DOCKER_COMPOSE_PATH) down

docker.build:
	docker build . -f $(DOCKER_PATH)

docker.compose.run:
	docker compose --env-file $(ENV_CONFIG) --file $(DOCKER_COMPOSE_PATH) up --build

docker.compose.up:
	docker compose --env-file $(ENV_CONFIG) --file $(DOCKER_COMPOSE_PATH) up --build --force-recreate -d

docker.compose.down:
	docker compose --env-file $(ENV_CONFIG) --file $(DOCKER_COMPOSE_PATH) down