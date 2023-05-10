ifeq ($(shell test -e '.env' && echo -n yes),yes)
	include .env
endif

# Commands
env:
	@$(eval SHELL:=/bin/bash)
	@cp .env.sample .env
	@echo "SECRET_KEY=$$(openssl rand -hex 32)" >> .env

build:
	go build -o ./cmd/main

run:
	go run ./cmd/main

unit-test:
	go test ./...

swagger:
	swag init -g cmd/main/main.go

df:
	docker build --tag simple_web_games .

service_up:
	docker-compose -f docker-compose.yml up -d --remove-orphans

.PHONY: env build run swagger df service_up unit-test