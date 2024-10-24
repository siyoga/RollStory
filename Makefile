.PHONY: build

build:
	go build -race -v cmd/main.go

build_docker:
	docker build --build-arg MODULE=bot --build-arg MODE=local -t rollstorybot . --no-cache

run:
	docker compose up -d --build

stop:
	docker compose down
