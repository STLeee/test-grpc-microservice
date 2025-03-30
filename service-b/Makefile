ENV=dev
VERSION := 1.0.0

build:
	docker build -t service-b:latest -f Dockerfile ../
	docker tag service-b:latest service-b:$(ENV)
	docker tag service-b:latest service-b:$(VERSION)
