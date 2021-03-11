run:
	go run ./...

build:
	go build ./...

docker-build:
	docker build -t jonathanwthom/meminders:latest .

docker-run:
	docker run --rm -v $(shell pwd)/meminders-dev.db:/app/meminders-dev.db jonathanwthom/meminders

lint:
	golangci-lint run

test:
	go test -v -cover ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''
