.PHONY: run build test test-cover docker-up docker-down migrate seed swagger clean

APP_NAME=dating-app

run:
	go run cmd/api/main.go

build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

test:
	go test -v -race ./internal/...

test-cover:
	go test -v -cover -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

migrate:
	go run cmd/api/main.go -migrate

seed:
	go run cmd/api/main.go -seed

swagger:
	~/go/bin/swag init -g main.go

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

clean:
	rm -rf bin coverage.out coverage.html docs
