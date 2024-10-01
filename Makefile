BUILD_FOLDER := build
BUILD_FILE := api

GO := go
GOMOD := $(GO) mod
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test

# .env should exist by the moment we start running make scripts
include .env
export

.PHONY: all build run clean tidy test up down psql

all: build run

build:
	$(GOBUILD) -race -o "$(BUILD_FOLDER)/$(BUILD_FILE)" -v ./cmd/api/main.go

run:
	"$(BUILD_FOLDER)/$(BUILD_FILE)"

# Executing `go clean` removes any non-make-related builds, generated using `go build`, such as the api.exe binary produced by the api package.
# `sudo rm -rf pgdata` will not return errors (related to rm at least)
clean: down
	$(GOCLEAN)
	-rm -r $(BUILD_FOLDER)
	sudo rm -rf pgdata/

tidy:
	$(GOMOD) tidy -e

test:
	$(GOTEST) ./tests -v

psql:
	docker exec -it $(DEBUG_POSTGRES_CONTAINER_NAME) psql -U $(DEBUG_POSTGRES_USER) -d $(DEBUG_POSTGRES_DATABASE)

up:
	docker compose up -d

down:
	docker compose down

# Generating mocks

SRC_FILES := $(wildcard internal/service/repository/*.go)
MOCK_FILES := $(patsubst internal/service/repository/%.go,internal/service/repository/mock_%.go,$(filter-out internal/service/repository/mock_%.go,$(SRC_FILES)))

mock: $(MOCK_FILES)

internal/service/repository/mock_%.go: internal/service/repository/%.go
	@if [ "$(@F)" != "mock_$(*F)" ]; then \
		mockgen -source="$<" -destination="$@" -package=repository; \
	fi

# Generating users

USERS := \
	'{"firstname": "John", "lastname": "Doe", "username": "johndoe", "email": "johndoe@example.com", "password": "T3st1ngP4$$w0rd"}', \
	'{"firstname": "Jane", "lastname": "Smith", "username": "janesmith", "email": "janesmith@example.com", "password": "P@ssw0rd123"}', \
	'{"firstname": "Alice", "lastname": "Johnson", "username": "alicejohnson", "email": "alicejohnson@example.com", "password": "qwerty"}', \
	'{"firstname": "Bob", "lastname": "Brown", "username": "bobbrown", "email": "bobbrown@example.com", "password": "password123"}'
ENDPOINT=http://localhost:8080/v1/auth/register

create_users:
	@echo "Creating users..."
	@for user_data in $(USERS); do \
		curl --silent --request POST \
			--url $(ENDPOINT) \
			--header 'Content-Type: application/json' \
			--data "$$user_data"; \
	done
