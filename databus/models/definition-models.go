// models.go
package models

import (
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* The state object for JSON, defined in definitions.json */
type StateJs struct {
	Hex   string `bson:"Hex" json:"Hex"`
	Label string `bson:"Label" json:"Label"`
}

/* The state object for database and internal use */
type StateRaw struct {
	Hex   uint16 `bson:"Hex" json:"Hex"`
	Label string `bson:"Label" json:"Label"`
}

/* The state object for the API */
// type StateDTO struct {
// 	Hex   string `bson:"Hex" json:"Hex"`
// 	Label string `bson:"Label" json:"Label"`
// }

/* The definition object for JSON, defined in definitions.json */
type DefinitionJs struct {
	Name        string    `bson:"Name" json:"Name"`
	Description string    `bson:"Description,omitempty" json:"Description,omitempty"`
	States      []StateJs `bson:"States" json:"States"`
}

/* The definition object for database and internal use */
type DefinitionRaw struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"ID,omitempty"`
	Name        string             `bson:"Name" json:"Name"`
	Description string             `bson:"Description,omitempty" json:"Description,omitempty"`
	States      []StateRaw         `bson:"States" json:"States"`
}

// --------------------- Conversion functions ---------------------
/*
Basic conversion order:
	Js -> Raw -> DTO
	DTO -> Raw -> Js
*/

func (m *DefinitionJs) ToRaw() DefinitionRaw {
	// Convert states
	states := make([]StateRaw, len(m.States))
	for i, st := range m.States {

		// Parse the hex string (e.g. "0x00", "0x01", etc.)
		val, err := strconv.ParseUint(st.Hex, 0, 16)
		if err != nil {
			log.Fatalf("invalid hex value: %q, error: %v", st.Hex, err)
		}
		states[i] = StateRaw{
			Hex:   uint16(val),
			Label: st.Label,
		}
	}

	return DefinitionRaw{
		Name:        m.Name,
		Description: m.Description,
		States:      states,
	}
}

func (m *DefinitionRaw) ToJs() DefinitionJs {
	// Convert states
	states := make([]StateJs, len(m.States))
	for i, st := range m.States {
		states[i] = StateJs{
			Hex:   fmt.Sprintf("%#02x", st.Hex),
			Label: st.Label,
		}
	}

	return DefinitionJs{
		Name:        m.Name,
		Description: m.Description,
		States:      states,
	}
}

// --------------------- Print functions ---------------------

func (a *DefinitionJs) Print() {
	fmt.Println("AppModel: ")
	fmt.Println("Name: " + a.Name)
	fmt.Println("Description: " + a.Description)
	fmt.Println("States: ")
	for i := 0; i < len(a.States); i++ {
		fmt.Println("\tHex: " + fmt.Sprintf("%#02x", a.States[i].Hex))
		fmt.Println("\tLabel: " + a.States[i].Label + "\n")
	}
	fmt.Println("------------")
}
