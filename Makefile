include .env
export

export PROJECT_ROOT=$(shell pwd)
COMPOSE_FILE=docker-compose.yaml

up:
	@docker compose -f $(COMPOSE_FILE) up -d --build

down:
	@docker compose -f $(COMPOSE_FILE) down

seed:
	@docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f /seeds/seed.sql

logs:
	@docker compose logs -f app

.PHONY: up down seed logs
