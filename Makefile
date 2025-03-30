ENV=dev

build-proto:
	protoc --proto_path=./proto --go_out=. --go-grpc_out=. ./proto/*.proto

build:
	$(MAKE) -C service-a build
	$(MAKE) -C service-b build

run:
	docker-compose -f docker-compose.$(ENV).yml up -d

stop:
	docker-compose -f docker-compose.$(ENV).yml down
