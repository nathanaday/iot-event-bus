# IoT Event Bus

## Disclaimers

>[!note]
> This repository contains proof of concepts and design pattern demonstrations. It is not an "out of the box" app that will automatically work with all data schemas and persistence layers. Feel free to run the examples and modify as needed :)

## Description

I still use this project as a reference or base files for building lightweight event busses. This project uses a Go application (`databus/`) to interface with a MongDB instance (docker hosted) and announce all entity changes on MQTT (eclipse-mosquitto).

**To break it down further:**

- Use MongoDB to store entity models, rules, and records. 
  - This can be anything with a schema and state (smart lights, sprinklers, sensors...)
- Use the Go-hosted REST API to view or modify entities 
  - (e.g. add a light, view light details, toggle lighting state)
- All entity state changes are emitted by (1) group membership or (2) individual entity using MQTT.
  - e.g. subscribe to `events/{entity_id}`  to catch a message like `{... "entity_id": 0x00, "new_state": 0}`
- Now you have a working event bus you can use to control decently complex population (hundreds, thousands) of sensors, controllers, and devices!

## What you can do next...

**For inspiration, you can use easily adapt this design pattern to:**

- Build a web-based dashboard to toggle smart lights in real time
- Manage a fleet of micro-controllers from a central controller
- Connect local smart controller events to a cloud of your choice

## Service Architecture

- **MQTT Broker (Mosquitto)**: Message broker on ports 1883 (MQTT) and 9001 (WebSockets)
- **MongoDB**: NoSQL database for persisting device data on port 27017
- **Go API Server**: REST API built with Gin framework on port 8080

## Quick Start

### Using Make (Easiest)

```bash
make up-build    # Start all services
make logs        # View logs
make test-api    # Test API endpoints
make test-mqtt   # Test MQTT connectivity
make down        # Stop all services
make help        # See all available commands
```

### Using Docker Compose (Manual)

```bash
cd use_mqtt
docker-compose up --build
```

This will start:
- MQTT Broker (ports 1883, 9001)
- MongoDB (port 27017)
- Go API Server (port 8080)


### Environment Variables

The application supports the following environment variables:

- `MQTT_BROKER_URL`: MQTT broker address (default: `tcp://localhost:1883`)
- `MONGODB_URI`: MongoDB connection string (default: `mongodb://localhost:27017`)
- `SERVER_ADDRESS`: Server bind address (default: `127.0.0.1:8080`)
- `DOCUMENTS_PATH`: Path to configuration JSON files (default: `/documents` in Docker, auto-detected locally)

### Have fun!



