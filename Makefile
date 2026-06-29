include .env
export

export PROJECT_ROOT=$(shell pwd)
COMPOSE_FILE=docker-compose.yaml

DB_SERVICE=postgres-container
APP_SERVICE=app-container

up:
	@docker compose -f $(COMPOSE_FILE) up -d --build

down:
	@docker compose -f $(COMPOSE_FILE) down

migrate-create:
	@if [ -z "$(seq)" ]; then \
    		  echo "seq variable is undefined"; \
    		  exit 1; \
    fi; \
    docker compose run --rm migration create \
    	-ext sql \
    	-dir /migrations \
    	-seq $(seq)


migrate-up:
	@if [ -z "$(num)" ]; then \
        		  echo "num variable is undefined"; \
        		  exit 1; \
    fi; \
	make migrate-action action=up num=$(num)

migrate-down:
	@if [ -z "$(num)" ]; then \
    	echo "num variable is undefined"; \
        exit 1; \
    fi; \
	@make migrate-action action=down num=$(num)

migrate-action:
	@if [ -z "$(action)" ] || [ -z "$(num)" ]; then \
  		echo "num or action variable is undefined"; \
        exit 1; \
    fi; \
    docker compose run --rm migration \
    -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
    -path /migrations \
    $(action) $(num)


.PHONY: up down

