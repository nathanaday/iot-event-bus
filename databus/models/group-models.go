package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* The group object for JSON, defined from groups.json */
type GroupJs struct {
	Name               string   `bson:"Name" json:"Name"`
	Description        string   `bson:"Description,omitempty" json:"Description,omitempty"`
	AllowedDefinitions []string `bson:"AllowedDefinitions" json:"AllowedDefinitions"`
}

/* For database and internal use */
type GroupRaw struct {
	ID                 primitive.ObjectID   `bson:"_id,omitempty" json:"ID,omitempty"`
	Name               string               `bson:"Name" json:"Name"`
	Description        string               `bson:"Description,omitempty" json:"Description,omitempty"`
	AllowedDefinitions []primitive.ObjectID `bson:"AllowedDefinitions" json:"AllowedDefinitions"`
}

// --------------------- Conversion functions ---------------------
/*
Basic conversion order:
	Js -> Raw -> DTO
	DTO -> Raw -> Js
*/

func (g *GroupJs) ToRaw(defs []DefinitionRaw) *GroupRaw {
	/*
		Converts a GroupJs object to a GroupRaw object
	*/

	defMap := make(map[string]primitive.ObjectID)
	for _, def := range defs {
		defMap[def.Name] = def.ID
	}

	raw := &GroupRaw{
		Name:               g.Name,
		Description:        g.Description,
		AllowedDefinitions: make([]primitive.ObjectID, len(g.AllowedDefinitions)),
	}

	fmt.Println("Creating space for ", g.AllowedDefinitions, " allowed definitions for group ", g.Name)

	for i, defName := range g.AllowedDefinitions {
		if id, exists := defMap[defName]; exists {
			raw.AllowedDefinitions[i] = id
		} else {
			// Handle missing definition (log, skip, or return error)
			fmt.Printf("Warning: definition '%s' not found\n", defName)
		}
	}
	return raw
}

func (g *GroupRaw) ToJs(definitions []DefinitionRaw) *GroupJs {
	/*
		Converts a GroupRaw object to a GroupJs object
	*/
	idToNameMap := make(map[primitive.ObjectID]string)
	for _, def := range definitions {
		idToNameMap[def.ID] = def.Name
	}

	// Initialize the Js group object
	js := &GroupJs{
		Name:               g.Name,
		Description:        g.Description,
		AllowedDefinitions: make([]string, len(g.AllowedDefinitions)),
	}

	// Map ObjectIDs back to names
	for i, objID := range g.AllowedDefinitions {
		if name, exists := idToNameMap[objID]; exists {
			js.AllowedDefinitions[i] = name
		} else {
			// Handle missing ObjectID gracefully
			fmt.Printf("Warning: ObjectID '%s' not found in definitions\n", objID.Hex())
		}
	}
	return js
}

// --------------------- Print functions ---------------------

func (g *GroupJs) Print() {
	fmt.Println("Group: ")
	fmt.Println("Name: " + g.Name)
	fmt.Println("Description: " + g.Description)
	fmt.Println("Allowed Models: ")
	for i := 0; i < len(g.AllowedDefinitions); i++ {
		fmt.Println(g.AllowedDefinitions[i])
	}
}
