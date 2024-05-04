BUILD_FOLDER=build
BUILD_FILE=api

GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test

include .env
export

all: build run

build: cmd/api/main.go
	$(GOBUILD) -race -o "$(BUILD_FOLDER)/$(BUILD_FILE)" -v $<

run:
	"$(BUILD_FOLDER)/$(BUILD_FILE)"

clean: down
	@# Use "$(GOCLEAN)" so it removes any self-compiled bins (not using `make build`)
	$(GOCLEAN)
	-rm -r $(BUILD_FOLDER)
	@# The rm should be ran by a superuser. The pgdata is created by docker using sudo
	-sudo rm -rf pgdata/

tidy:
	go mod tidy -e

up:
	docker compose up -d

down:
	docker compose down

psql:
	docker exec -it arkhon-db psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

.PHONY: build run all
