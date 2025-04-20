ENV=dev

build-proto:
	protoc --proto_path=./proto --go_out=. --go-grpc_out=. ./proto/*.proto

build:
	$(MAKE) -C service-api build
	$(MAKE) -C service-a build
	$(MAKE) -C service-b build

test:
	$(MAKE) -C service-api test
	$(MAKE) -C service-a test
	$(MAKE) -C service-b test

run:
	docker-compose -f docker-compose.$(ENV).yml up -d

stop:
	docker-compose -f docker-compose.$(ENV).yml down
