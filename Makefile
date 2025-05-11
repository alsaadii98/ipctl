# Variables
GO = go
BUILD_DIR = ./bin
BINARY_NAME = ipctl
PORT ?= 8080


.PHONY: all
all: build


.PHONY: build
build:
	@echo "Building the application..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .


.PHONY: run
run: build
	@echo "Running the application on port $(PORT)..."
	$(BUILD_DIR)/$(BINARY_NAME)


.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf $(BUILD_DIR)


.PHONY: install
install:
	@echo "Installing Go modules..."
	$(GO) mod tidy
	$(GO) mod vendor