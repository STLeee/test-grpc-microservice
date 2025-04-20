APP_NAME=service-b

ENV=dev
VERSION=1.0.0

run:
	go run main.go

build:
	docker build -t $(APP_NAME):latest -f Dockerfile ../
	docker tag $(APP_NAME):latest $(APP_NAME):$(ENV)
	docker tag $(APP_NAME):latest $(APP_NAME):$(VERSION)

make test:
	go fmt
	go vet ./...
	go test ./... -coverprofile=coverage.out
