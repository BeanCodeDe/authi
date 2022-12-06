SRC_PATH?=./cmd/authi
APP_NAME?=authi
DOCKER_COMPOSE_PATH?=./deployments/docker-compose-postgres.yml
DOCKER_PATH?=./build/Dockerfile

version.up:
	bash ./scripts/auto-increment-version.sh

init.token:
	sh ./scripts/generateKeyFile.sh

app.build:
	go mod download
	go build -o $(APP_NAME) $(SRC_PATH)

app.ut.run:
	go test ./internal/... -v

app.jt.run:
	docker compose --file $(DOCKER_COMPOSE_PATH) up --build --force-recreate -d
	go test ./test
	docker compose --file $(DOCKER_COMPOSE_PATH) down

docker.build:
	docker build . -f $(DOCKER_PATH)

docker.compose.run:
	docker compose --file $(DOCKER_COMPOSE_PATH) up --build
