package handlers

import (
	"databus/models"
	"databus/persistence"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllDefinitionsHandler(g *gin.Context) {

	dfs, err := persistence.GetAllDefinitions()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Convert the fetched models to their js/DTO representations
	dfs_dto := make([]models.DefinitionJs, len(dfs))
	for i := range dfs {
		dfs_dto[i] = dfs[i].ToJs()
	}

	g.JSON(200, dfs_dto)
}

func GetAllGroupsHandler(g *gin.Context) {

	groups, err := persistence.GetAllGroups()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	dfRaw, err := persistence.GetAllDefinitions()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	gps_dto := make([]models.GroupJs, len(groups))
	for i := range groups {
		gps_dto[i] = *groups[i].ToJs(dfRaw)
	}

	g.JSON(200, gps_dto)
}

func GetAllReactiveEntitiesHandler(g *gin.Context) {

	reactiveEntities, err := persistence.GetAllReactiveEntities()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	definitions, err := persistence.GetAllDefinitions()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	groups, err := persistence.GetAllGroups()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	entities_dto := make([]models.ReactiveEntityJs, len(reactiveEntities))
	for i := range reactiveEntities {
		entities_dto[i] = *reactiveEntities[i].ToJs(definitions, groups)
	}

	g.JSON(200, entities_dto)
}

func GetDefinitionByNameHandler(g *gin.Context) {

	name := g.Param("definitionName")

	def, err := persistence.GetDefinitionByName(name)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, def.ToJs())
}

func GetGroupByNameHandler(g *gin.Context) {

	name := g.Param("groupName")

	group, err := persistence.GetGroupByName(name)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	dfRaw, err := persistence.GetAllDefinitions()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, group.ToJs(dfRaw))
}

func GetReactiveEntityByHexHandler(g *gin.Context) {

	hex := g.Param("entityHex") // string

	// string -> uint16
	hexInt64, err := strconv.ParseUint(hex, 16, 16)
	if err != nil {
		g.JSON(400, gin.H{"error": err.Error()})
		return
	}
	hexInt := uint16(hexInt64)

	reactiveEntity, err := persistence.GetReactiveEntityByHex(hexInt)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	definitions, err := persistence.GetAllDefinitions()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	groups, err := persistence.GetAllGroups()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	g.JSON(200, reactiveEntity.ToJs(definitions, groups))
}

func GetReactiveEntitiesByGroupHandler(g *gin.Context) {

	gl := g.Param("groupList")
	groupNames := strings.Split(gl, ",")

	reactiveEntities, err := persistence.GetReactiveEntitiesByGroup(groupNames)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	definitions, err := persistence.GetAllDefinitions()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	groups, err := persistence.GetAllGroups()
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}

	reactiveEntitiesJs := make([]models.ReactiveEntityJs, len(reactiveEntities))
	for i := range reactiveEntities {
		reactiveEntitiesJs[i] = *reactiveEntities[i].ToJs(definitions, groups)
	}

	g.JSON(200, reactiveEntitiesJs)
}
