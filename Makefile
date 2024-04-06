BUILD_FOLDER=build
BUILD_FILE=api

GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test

all: build run

build: cmd/api/main.go
	$(GOBUILD) -race -o "$(BUILD_FOLDER)/$(BUILD_FILE)" -v $<

run:
	PORT=8080 "$(BUILD_FOLDER)/$(BUILD_FILE)"

clean:
	@# Use "$(GOCLEAN)" so it removes any self-compiled bins (not using `make build`)
	$(GOCLEAN)
	-rm -r $(BUILD_FOLDER)

.PHONY: build run all