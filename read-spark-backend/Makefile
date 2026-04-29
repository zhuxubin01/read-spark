.PHONY: build run test migrate-up migrate-down

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/readspark?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/readspark?sslmode=disable" down
