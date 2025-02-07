# Define variables
DOCKER_COMPOSE = docker-compose
GO = go
APP_NAME = main.go

# Targets
.PHONY: all up down build run clean

# Default target to start the application
all: up run

# Bring up the MongoDB container using Docker Compose
up:
	$(DOCKER_COMPOSE) up -d

# Stop and remove MongoDB container
down:
	$(DOCKER_COMPOSE) down

# Build the Go application
build:
	$(GO) build -o app $(APP_NAME)

# Run the Go application
run: up
	sleep 10  # Wait for MongoDB to initialize
	$(GO) run $(APP_NAME)

# Clean up generated files
clean:
	rm -f app
	$(DOCKER_COMPOSE) down
