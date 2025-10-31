#!/bin/bash

# Simple script to test MQTT connectivity
# Usage: ./test_mqtt.sh

echo "Testing MQTT connection..."

# Publish a test message
docker exec mqtt5 mosquitto_pub -h localhost -t "test/topic" -m "Hello from Databus Backend" -q 1

echo "Test message published to topic: test/topic"
echo ""
echo "To subscribe to messages, run:"
echo "docker exec -it mqtt5 mosquitto_sub -h localhost -t '#' -v"

