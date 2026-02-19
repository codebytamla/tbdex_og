SHELL := /bin/sh

GOOSE ?= goose
GOOSE_DIR ?= migrations
SEED_FILE ?= seeds/001_roles_permissions.sql

.PHONY: run migrate-up migrate-down seed seed-admin setup-db docs-lint docs-build docs-preview docs-ui

run:
	go run ./cmd/server

migrate-up:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL is required"; exit 1; fi
	$(GOOSE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" up

migrate-down:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL is required"; exit 1; fi
	$(GOOSE) -dir $(GOOSE_DIR) postgres "$$DATABASE_URL" down

seed:
	@if [ -z "$$DATABASE_URL" ]; then echo "DATABASE_URL is required"; exit 1; fi
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f $(SEED_FILE)

seed-admin:
	@if [ -z "$$ADMIN_EMAIL" ]; then echo "ADMIN_EMAIL is required"; exit 1; fi
	@if [ -z "$$ADMIN_PASSWORD" ]; then echo "ADMIN_PASSWORD is required"; exit 1; fi
	go run ./cmd/seed

setup-db: migrate-up seed

