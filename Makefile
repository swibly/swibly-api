BUILD_FOLDER := build
BUILD_FILE := api

GO := go
GOMOD := $(GO) mod
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test

include .env
export

.PHONY: all build run clean tidy test up down psql

all: up

build: $(BUILD_FOLDER)/$(BUILD_FILE)

$(BUILD_FOLDER)/$(BUILD_FILE): ./cmd/api/main.go
	@mkdir -p $(BUILD_FOLDER)
	$(GOBUILD) -race -o "$@" -v $<

run: $(BUILD_FOLDER)/$(BUILD_FILE)
	@echo "Starting $(BUILD_FILE) from $(BUILD_FOLDER)..."
	"./$<"

clean: down
	$(GOCLEAN)
	@if [ -d "$(BUILD_FOLDER)" ]; then rm -r "$(BUILD_FOLDER)"; fi
	@if [ -d "pgdata" ]; then sudo rm -rf pgdata/; fi

tidy:
	$(GOMOD) tidy -e

test: TEST_DIR ?= ./tests
test:
	$(GOTEST) $(TEST_DIR) -v

psql:
	@docker exec -it swibly-api-db psql -U "$(POSTGRES_USER)" -d "$(POSTGRES_DATABASE)" || \
	echo "Failed to connect to database. Check container and environment variables."

up:
	@docker compose up -d postgres
	@docker compose up --build --no-deps swibly-api

down:
	@docker compose down
