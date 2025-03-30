VERSION := 1.0.0

build:
	docker build -t service-a:latest -f Dockerfile ../
	docker tag service-a:latest service-a:$(VERSION)
