// main.go
package main

import (
	"databus/cmd/api"
	"databus/cmd/config"
	"databus/network"
	"databus/persistence"
)

// var mongoClient *mongo.Client

func main() {

	// Initialize MQTT client
	network.InitMQTTClient()

	// Start the WebSocket client (listens for FCodes)
	// go startWebSocketClient()

	// Start the WebSocket server (sends notifications)
	// go startWebSocketServer()
	persistence.Connect()
	defer persistence.Disconnect()

	// Configuration parsing
	config.ParseAllConfigs()

	// Initialize the router and routes
	api.InitializeRoutes()
}
