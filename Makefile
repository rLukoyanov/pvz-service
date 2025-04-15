BINARY_NAME=pvz-service
DB_URL=postgres://postgres:postgres@localhost:5432/pvz?sslmode=disable
MIGRATIONS_DIR=./migrations

.PHONY: all build run clean docker-up docker-down migrate-up migrate-down generate

all: build

up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

up_build: build
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BINARY_NAME} ./cmd/server/main.go

run:
	go run ./cmd/main.go

clean:
	rm -f $(BINARY_NAME)