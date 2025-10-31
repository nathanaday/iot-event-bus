package config

import (
	"databus/models"
	"databus/persistence"
	"databus/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var path string

type documents struct {
	definitions string
	groups      string
}

var jsons documents

func init() {
	// Check for environment variable first (for Docker/production)
	if documentsPath := os.Getenv("DOCUMENTS_PATH"); documentsPath != "" {
		path = documentsPath
	} else {
		// Get the directory of the current file
		_, filename, _, _ := runtime.Caller(0)
		dir := filepath.Dir(filename)

		// Traverse up to the project root (one level up from the go.mod file)
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				dir = filepath.Dir(dir) // Move one level up from the go.mod directory
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				// If go.mod not found, try default path for Docker
				if _, err := os.Stat("/documents"); err == nil {
					path = "/documents"
					log.Println("Using default documents path: /documents")
					jsons = documents{
						definitions: "definitions.json",
						groups:      "groups.json",
					}
					return
				}
				log.Fatal("go.mod file not found and DOCUMENTS_PATH not set")
			}
			dir = parent
		}

		// Set the path to the documents directory
		path = filepath.Join(dir, "documents")
	}

	jsons = documents{
		definitions: "definitions.json",
		groups:      "groups.json",
	}
}

func ParseAllConfigs() {
	/*
		Order matters! The hierarchy for validation is designed like so:
		- Definitions are isolated objects that do not refer/link to any other config, so they can be parsed first
		- Groups refer to definitions, so validation considers whether the model exists before comitting them

		For each config type, the three steps are executed:
		1. Parse from JSON file
		2. Validate the parsed data based on rules
		3. Insert the validated data into the database

		Note: Reactive entities are managed via API and not loaded from static files.
	*/

	// --- --- --- --- --- --- Definitions --- --- --- --- --- ---
	fmt.Println("\nParsing and validating Definition configuration...")
	djs, err := ParseDefinitions()
	if err != nil {
		log.Fatal(utils.StrToRed("Error parsing definition json: "), err)
	}
	fmt.Println(utils.StrToGreen("\tLoaded definition json"))

	vms, err := ValidateDefinitions(djs)
	if err != nil {
		log.Fatal(utils.StrToRed("Error validating definitions: "), err)
	}
	fmt.Println(utils.StrToGreen("\tValidated definition json"))

	_, err = persistence.InsertDefinitions(persistence.MongoClient, vms)
	if err != nil {
		log.Fatal(utils.StrToRed("Error inserting definitions: "), err)
	}
	fmt.Println(utils.StrToGreen("\tInserted definitions into persistence"))

	// --- --- --- --- --- --- Groups --- --- --- --- --- ---
	fmt.Println("\nParsing and validating Groups configuration...")
	gps, err := ParseGroups()
	if err != nil {
		log.Fatal(utils.StrToRed("Error parsing groups.json: "), err)
	}
	fmt.Println(utils.StrToGreen("\tLoaded groups.json"))
	vgps, err := ValidateGroups(gps, vms)
	if err != nil {
		log.Fatal(utils.StrToRed("Error validating groups: "), err)
	}
	fmt.Println(utils.StrToGreen("\tValidated groups.json"))
	_, err = persistence.InsertGroups(persistence.MongoClient, vgps)
	if err != nil {
		log.Fatal(utils.StrToRed("Error inserting groups: "), err)
	}
	fmt.Println(utils.StrToGreen("\tInserted groups into persistence"))

	fmt.Println(utils.StrToGreen("\nAll configurations parsed, validated, and inserted successfully!\n"))
	fmt.Println(utils.StrToGreen("Note: Reactive entities are managed via API endpoints.\n"))

}

func ParseDefinitions() ([]models.DefinitionJs, error) {
	// Parse definition json

	jsf := filepath.Join(path, jsons.definitions)

	// Open the file
	file, err := os.Open(jsf)
	if err != nil {
		log.Fatal("Error opening groups.json:", err)
		return nil, err
	}
	defer file.Close()

	// Read the file
	byteValue, _ := io.ReadAll(file)

	var definitions []models.DefinitionJs

	json.Unmarshal(byteValue, &definitions)

	// for i := 0; i < len(definitions); i++ {
	// 	definitions[i].Print()
	// }

	return definitions, err
}

func ParseGroups() ([]models.GroupJs, error) {
	// Parse groups.json

	jsf := filepath.Join(path, jsons.groups)

	// Open the file
	file, err := os.Open(jsf)
	if err != nil {
		log.Fatal("Error opening groups.json:", err)
		return nil, err
	}
	defer file.Close()

	// Read the file
	byteValue, _ := io.ReadAll(file)

	var groups []models.GroupJs

	json.Unmarshal(byteValue, &groups)

	// for i := 0; i < len(groups); i++ {
	// 	groups[i].Print()
	// }

	return groups, err

}

