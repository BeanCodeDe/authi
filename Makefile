SRC_PATH?=./internal
APP_NAME?=auth

build:
	go build -buildvcs=false -o $(APP_NAME) $(SRC_PATH)