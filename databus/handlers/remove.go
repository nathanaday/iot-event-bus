package handlers

import (
	"databus/persistence"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteReactiveEntityHandler deletes a reactive entity by its hex ID
func DeleteReactiveEntityHandler(g *gin.Context) {
	hex := g.Param("entityHex") // string

	// string -> uint16
	hexInt64, err := strconv.ParseUint(hex, 16, 16)
	if err != nil {
		g.JSON(400, gin.H{"error": "Invalid hex format", "details": err.Error()})
		return
	}
	hexInt := uint16(hexInt64)

	// Delete the reactive entity
	deletedCount, err := persistence.DeleteReactiveEntityByHex(hexInt)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to delete reactive entity", "details": err.Error()})
		return
	}

	if deletedCount == 0 {
		g.JSON(404, gin.H{"error": "Reactive entity not found"})
		return
	}

	g.JSON(200, gin.H{"message": "Reactive entity deleted successfully", "entityHex": hex})
}
