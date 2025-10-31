// reactive-entity-models.go
package models

import (
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* The data object, default/empty data to be added to the reactive entity object on creation */
type DataObj struct {
	CurrentState int       `bson:"CurrentState" json:"CurrentState"`
	LastUpdated  time.Time `bson:"LastUpdated" json:"LastUpdated"`
}

/* The location object, for reactive entities */
type Location struct {
	SLCoordX float64 `bson:"SLCoordX" json:"SLCoordX"`
	SLCoordY float64 `bson:"SLCoordY" json:"SLCoordY"`
	Name     string  `bson:"Name" json:"Name"`
	Rack     int     `bson:"Rack" json:"Rack"`
}

/* The reactive entity object for JSON/API */
type ReactiveEntityJs struct {
	EntityHex   string   `bson:"EntityHex" json:"EntityHex"`
	Description string   `bson:"Description,omitempty" json:"Description,omitempty"`
	Location    Location `bson:"Location" json:"Location"`
	Definition  string   `bson:"Definition" json:"Definition"`
	Groups      []string `bson:"Groups" json:"Groups"`
	Data        DataObj  `bson:"Data" json:"Data"`
}

/* The reactive entity object for database and internal use */
type ReactiveEntityRaw struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"ID,omitempty"`
	EntityHex   uint16               `bson:"EntityHex" json:"EntityHex"`
	Description string               `bson:"Description,omitempty" json:"Description,omitempty"`
	Location    Location             `bson:"Location" json:"Location"`
	Definition  primitive.ObjectID   `bson:"Definition" json:"Definition"`
	Groups      []primitive.ObjectID `bson:"Groups" json:"Groups"`
	Data        DataObj              `bson:"Data" json:"Data"`
}

// --------------------- Conversion functions ---------------------
/*
Basic conversion order:
	Js -> Raw -> DTO
	DTO -> Raw -> Js
*/

func (e *ReactiveEntityJs) ToRaw(definitions []DefinitionRaw, groups []GroupRaw) *ReactiveEntityRaw {
	/*
	   Converts a ReactiveEntityJs object to a ReactiveEntityRaw object.
	   Uses the provided definitions and groups to map names to ObjectIDs.
	*/

	// Create lookup maps for efficient mapping
	defMap := make(map[string]primitive.ObjectID)
	for _, def := range definitions {
		defMap[def.Name] = def.ID
	}

	groupMap := make(map[string]primitive.ObjectID)
	for _, group := range groups {
		groupMap[group.Name] = group.ID
	}

	// Convert EntityHex from string to uint16
	val, err := strconv.ParseUint(e.EntityHex, 0, 16)
	if err != nil {
		fmt.Printf("Error parsing EntityHex: %v\n", err)
	}

	// Initialize the raw reactive entity
	raw := &ReactiveEntityRaw{
		EntityHex:   uint16(val),
		Description: e.Description,
		Location:    e.Location,
		Data:        DataObj{},
	}

	// Map the Definition field
	if id, exists := defMap[e.Definition]; exists {
		raw.Definition = id
	} else {
		fmt.Printf("Warning: Definition '%s' not found\n", e.Definition)
	}

	// Map the Groups field
	for _, groupName := range e.Groups {
		if id, exists := groupMap[groupName]; exists {
			raw.Groups = append(raw.Groups, id)
		} else {
			fmt.Printf("Warning: Group '%s' not found\n", groupName)
		}
	}

	return raw
}

func (e *ReactiveEntityRaw) ToJs(definitions []DefinitionRaw, groups []GroupRaw) *ReactiveEntityJs {
	/*
	   Converts a ReactiveEntityRaw object to a ReactiveEntityJs object.
	   Uses the provided definitions and groups to map ObjectIDs back to names.
	*/

	// Create reverse lookup maps
	idToDefNameMap := make(map[primitive.ObjectID]string)
	for _, def := range definitions {
		idToDefNameMap[def.ID] = def.Name
	}

	idToGroupNameMap := make(map[primitive.ObjectID]string)
	for _, group := range groups {
		idToGroupNameMap[group.ID] = group.Name
	}

	// Initialize the JS reactive entity
	js := &ReactiveEntityJs{
		EntityHex:   fmt.Sprintf("%#02x", e.EntityHex),
		Description: e.Description,
		Location:    e.Location,
		Groups:      make([]string, len(e.Groups)),
		Data:        e.Data,
	}

	// Map the Definition field
	if name, exists := idToDefNameMap[e.Definition]; exists {
		js.Definition = name
	} else {
		fmt.Printf("Warning: Definition ID '%s' not found\n", e.Definition)
	}

	// Map the Groups field
	for j, group := range e.Groups {
		if name, exists := idToGroupNameMap[group]; exists {
			js.Groups[j] = name
		} else {
			fmt.Printf("Warning: Group ID '%s' not found\n", group)
		}
	}

	return js
}

// --------------------- Print functions ---------------------

func (e *ReactiveEntityJs) Print() {
	fmt.Println("EntityHex: " + fmt.Sprintf("%#02x", e.EntityHex))
	fmt.Println("Description: " + e.Description)
	fmt.Println("Location: ")
	fmt.Println("\tName: " + e.Location.Name)
	fmt.Println("\tSLCoordX: " + fmt.Sprintf("%f", e.Location.SLCoordX))
	fmt.Println("\tSLCoordY: " + fmt.Sprintf("%f", e.Location.SLCoordY))
	fmt.Println("\tRack: " + fmt.Sprintf("%d", e.Location.Rack))
	fmt.Println("Definition: " + e.Definition)
	fmt.Println("Groups: ")
	for j := 0; j < len(e.Groups); j++ {
		fmt.Println(e.Groups[j])
	}
}

func (d *DataObj) Print() {
	fmt.Println("CurrentState: " + fmt.Sprintf("%d", d.CurrentState))
	fmt.Println("LastUpdated: " + d.LastUpdated.String())
}
