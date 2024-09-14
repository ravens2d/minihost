BINARY_NAME=minihost

.PHONY: build run clean test docker-build docker-run docker-stop docker-clean docker-dev create-db

build:
	go build -o $(BINARY_NAME) ./cmd/main

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	go clean

test:
	go test ./...

docker-build:
	docker compose build

docker-run:
	docker compose up

docker-stop:
	docker compose down

docker-clean:
	docker compose down --rmi all --volumes

docker-dev: docker-build docker-run

create-db:
	sqlite3 database.db < schema/database.sql
