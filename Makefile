BUILD_FOLDER := build
BUILD_FILE := api

GO := go
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test

# .env should exist by the moment we start running make scripts
include .env
export

.PHONY: all build run clean tidy up down psql

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
	go mod tidy -e

psql:
	docker exec -it arkhon-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

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
