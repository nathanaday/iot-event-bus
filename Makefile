.PHONY: help up down logs build rebuild clean test-mqtt test-api status

help: ## Show this help message
	@echo "Backend - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Start all services (MQTT, MongoDB, API)
	cd use_mqtt && docker-compose up -d

up-build: ## Start all services and rebuild
	cd use_mqtt && docker-compose up -d --build

down: ## Stop all services
	cd use_mqtt && docker-compose down

logs: ## Show logs from all services
	cd use_mqtt && docker-compose logs -f

logs-api: ## Show logs from API service only
	cd use_mqtt && docker-compose logs -f databus-api

logs-mqtt: ## Show logs from MQTT broker only
	cd use_mqtt && docker-compose logs -f mqtt5

logs-mongo: ## Show logs from MongoDB only
	cd use_mqtt && docker-compose logs -f mongodb

build: ## Build the Go API service
	cd use_mqtt && docker-compose build databus-api

rebuild: ## Rebuild and restart all services
	cd use_mqtt && docker-compose down && docker-compose up -d --build

clean: ## Stop services and remove volumes
	cd use_mqtt && docker-compose down -v

status: ## Show status of all services
	cd use_mqtt && docker-compose ps

test-mqtt: ## Test MQTT connectivity
	@echo "Publishing test message to MQTT broker..."
	docker exec mqtt5 mosquitto_pub -h localhost -t "test/databus" -m "Test message from Makefile" -q 1
	@echo "✓ Message published successfully!"
	@echo ""
	@echo "To subscribe and see messages, run:"
	@echo "  docker exec -it mqtt5 mosquitto_sub -h localhost -t '#' -v"

test-api: ## Test API endpoints
	@echo "Testing API endpoints..."
	@echo ""
	@echo "GET /api/definitions:"
	@curl -s http://localhost:8080/api/definitions | jq . || echo "API not responding or jq not installed"
	@echo ""
	@echo "GET /api/groups:"
	@curl -s http://localhost:8080/api/groups | jq . || echo "API not responding or jq not installed"
	@echo ""
	@echo "GET /api/reactive-entities:"
	@curl -s http://localhost:8080/api/reactive-entities | jq . || echo "API not responding or jq not installed"

dev: ## Run API locally (requires local Go installation)
	cd databus && go run cmd/entrypoint/main.go

dev-deps: ## Start only MQTT and MongoDB for local development
	cd use_mqtt && docker-compose up -d mqtt5 mongodb

