APP_NAME=service-api

ENV=dev
VERSION=1.0.0

run:
	swag init
	go run main.go

build:
	docker build -t $(APP_NAME):latest -f Dockerfile ../
	docker tag $(APP_NAME):latest $(APP_NAME):$(ENV)
	docker tag $(APP_NAME):latest $(APP_NAME):$(VERSION)
