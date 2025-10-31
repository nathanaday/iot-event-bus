package api

import (
	"databus/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes() {
	router := gin.Default()

	router.SetTrustedProxies(nil)
	// router.SetTrustedProxies([]string{"192.168.1.2"})  // Example

	// ------------ API Endpoints ------------
	// Definitions API
	router.GET("/api/definitions", handlers.GetAllDefinitionsHandler)
	router.GET("/api/definitions/:definitionName", handlers.GetDefinitionByNameHandler)

	// Groups API
	router.GET("/api/groups", handlers.GetAllGroupsHandler)
	router.GET("/api/groups/:groupName", handlers.GetGroupByNameHandler)

	// Reactive Entities API
	router.GET("/api/reactive-entities", handlers.GetAllReactiveEntitiesHandler)
	router.GET("/api/reactive-entities/byHex/:entityHex", handlers.GetReactiveEntityByHexHandler)
	router.GET("/api/reactive-entities/byGroups/:groupList", handlers.GetReactiveEntitiesByGroupHandler)
	router.POST("/api/reactive-entities", handlers.CreateReactiveEntityHandler)
	router.DELETE("/api/reactive-entities/:entityHex", handlers.DeleteReactiveEntityHandler)
	// TODO: Add PUT endpoint for updating reactive entities
	// router.PUT("/api/reactive-entities/:entityHex", handlers.UpdateReactiveEntityHandler)

	// ------------ Groups API ------------
	// router.GET("/groups", handlers.GetGroupsHandler)

	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "127.0.0.1:8080"
	}
	router.Run(serverAddr) // Start the server
}
