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


.PHONY: env build run