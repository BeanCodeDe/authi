SRC_PATH?=./cmd/authi
APP_NAME?=authi
DOCKER_COMPOSE_PATH?=./deployments/docker-compose.yml
DOCKER_PATH?=./build/Dockerfile
ENV_CONFIG?=./deployments/dev.env

app.build:
	go build -o $(APP_NAME) $(SRC_PATH)

docker.build:
	docker build . -f $(DOCKER_PATH)

docker.cpmpose.run:
	docker compose --env-file $(ENV_CONFIG) --file $(DOCKER_COMPOSE_PATH) up --build --force-recreate