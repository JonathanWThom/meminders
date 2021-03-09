run:
	go run ./...

docker-build:
	docker build -t jonathanwthom/meminders:latest .

docker-run:
	docker run --rm -v $(shell pwd)/meminders-dev.db:/app/meminders-dev.db jonathanwthom/meminders

lint:
	golangci-lint run

test:
	go test -v -cover ./...
