# yq needed for parsing toml config files
# https://https://github.com/mikefarah/yq

# Makefile for running tests with different environments
# ENV can be set to DEVELOPMENT or PRODUCTION
ENV := DEVELOPMENT

ifeq ($(ENV),DEVELOPMENT)
CONFIG_PATH := ./internal/secret/config.dev.toml
else
CONFIG_PATH := ./internal/secret/config.prod.toml
endif

DB_URL := $(shell yq e '.database.url' $(CONFIG_PATH))
MIGRATE_DIR := ./migrations

MIGRATE := migrate -path $(MIGRATE_DIR) -database "$(DB_URL)"

.PHONY: migrate-up migrate-down migrate-goto migrate-force migrate-create

migrate-up:
	@$(MIGRATE) up

migrate-down:
	@$(MIGRATE) down 1

migrate-goto:
	@if [ -z "$(v)" ]; then \
		echo "error: input specific version. ex: make migrate-goto v=1"; \
		exit 1; \
	fi
	@$(MIGRATE) goto $(v)

migrate-clean:
	@if [ -z "$(v)" ]; then \
		echo "error: input specific version. ex: make migrate-clean v=1"; \
		exit 1; \
	fi
	@$(MIGRATE) force $(v)

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "error: input migration name. ex: make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)