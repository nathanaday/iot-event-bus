package handlers

import (
	"databus/models"
	"databus/persistence"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateReactiveEntityHandler creates a new reactive entity
func CreateReactiveEntityHandler(g *gin.Context) {
	var reactiveEntityJs models.ReactiveEntityJs

	// Bind JSON to the model
	if err := g.ShouldBindJSON(&reactiveEntityJs); err != nil {
		g.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validate required fields
	if reactiveEntityJs.EntityHex == "" {
		g.JSON(400, gin.H{"error": "EntityHex is required"})
		return
	}
	if reactiveEntityJs.Definition == "" {
		g.JSON(400, gin.H{"error": "Definition is required"})
		return
	}

	// Check if entity with this hex already exists
	existingEntity, _ := persistence.GetReactiveEntityByHex(parseEntityHex(reactiveEntityJs.EntityHex))
	if existingEntity != nil {
		g.JSON(409, gin.H{"error": "Reactive entity with this EntityHex already exists"})
		return
	}

	// Fetch definitions and groups for conversion
	definitions, err := persistence.GetAllDefinitions()
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to fetch definitions", "details": err.Error()})
		return
	}

	groups, err := persistence.GetAllGroups()
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to fetch groups", "details": err.Error()})
		return
	}

	// Validate that the definition exists
	definitionExists := false
	for _, def := range definitions {
		if def.Name == reactiveEntityJs.Definition {
			definitionExists = true
			break
		}
	}
	if !definitionExists {
		g.JSON(400, gin.H{"error": fmt.Sprintf("Definition '%s' not found", reactiveEntityJs.Definition)})
		return
	}

	// Validate that all groups exist
	groupMap := make(map[string]primitive.ObjectID)
	for _, group := range groups {
		groupMap[group.Name] = group.ID
	}
	for _, groupName := range reactiveEntityJs.Groups {
		if _, exists := groupMap[groupName]; !exists {
			g.JSON(400, gin.H{"error": fmt.Sprintf("Group '%s' not found", groupName)})
			return
		}
	}

	// Convert to Raw format
	reactiveEntityRaw := reactiveEntityJs.ToRaw(definitions, groups)

	// Insert into database
	_, err = persistence.InsertReactiveEntity(reactiveEntityRaw)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to create reactive entity", "details": err.Error()})
		return
	}

	// Convert back to Js for response
	createdEntity := reactiveEntityRaw.ToJs(definitions, groups)

	g.JSON(201, gin.H{
		"message": "Reactive entity created successfully",
		"entity":  createdEntity,
	})
}

// Helper function to parse entity hex
func parseEntityHex(hexStr string) uint16 {
	// Remove "0x" prefix if present
	if len(hexStr) > 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}
	
	var result uint16
	fmt.Sscanf(hexStr, "%x", &result)
	return result
}
