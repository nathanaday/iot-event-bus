package config

import (
	"databus/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateDefinitions(definitions []models.DefinitionJs) ([]models.DefinitionRaw, error) {
	// Validate the definitions
	// - verify all 'name' fields are unique
	// - verify states field is not empty
	// - verify states field contains unique states

	// Map to check for unique 'Name' fields across definitions
	nameMap := make(map[string]struct{})

	// Iterate through all definitions
	for _, df := range definitions {
		// Check for duplicate 'Name' fields
		if _, exists := nameMap[df.Name]; exists {
			return nil, fmt.Errorf("duplicate definition name detected: %s", df.Name)
		}
		nameMap[df.Name] = struct{}{}

		// Check that 'States' field is not empty
		if len(df.States) == 0 {
			return nil, fmt.Errorf("definition '%s' has an empty states list", df.Name)
		}

		// Verify uniqueness of 'Hex' in States using a map
		stateHexMap := make(map[string]struct{})
		for _, state := range df.States {
			if _, exists := stateHexMap[state.Hex]; exists {
				return nil, fmt.Errorf(
					"duplicate state hex value '%#02x' detected in definition '%s'",
					state.Hex, df.Name,
				)
			}
			stateHexMap[state.Hex] = struct{}{}
		}
	}

	valid_definitions := make([]models.DefinitionRaw, len(definitions))
	for i, df := range definitions {
		valid_definitions[i] = df.ToRaw()
	}

	// If all validations pass, return the definitions
	return valid_definitions, nil
}

func ValidateGroups(groups []models.GroupJs, validDefinitions []models.DefinitionRaw) ([]models.GroupRaw, error) {
	// Validate the groups
	// - verify all 'name' fields are unique
	// - verify definition tag referenced in 'AllowedModels' exists in definitions

	// Map to track unique group names
	groupNameMap := make(map[string]struct{})

	// Map to contain validated groups (GroupExt)
	var validGroups []models.GroupRaw

	// Map to track valid definion names for fast lookup
	DefNameMap := make(map[string]primitive.ObjectID)
	for _, def := range validDefinitions {
		DefNameMap[def.Name] = def.ID
	}

	// Validate groups
	for _, group := range groups {

		// Check for duplicate group names
		if _, exists := groupNameMap[group.Name]; exists {
			return nil, fmt.Errorf("duplicate group name detected: %s", group.Name)
		}
		groupNameMap[group.Name] = struct{}{}

		// Check that all referenced definition names in 'AllowedDefinitions' exist
		for _, ad := range group.AllowedDefinitions {
			if _, exists := DefNameMap[ad]; !exists {
				return nil, fmt.Errorf(
					"invalid definition reference '%s' in group '%s', not found in definitions",
					ad, group.Name,
				)
			}
		}

		// By this point the group is valid, so we can add it to the list of valid groups
		// with ObjectID's instead of names
		vg := group.ToRaw(validDefinitions)
		validGroups = append(validGroups, *vg)
	}

	// 	valGroup := models.GroupRaw{
	// 		ID:            group.ID,
	// 		Name:          group.Name,
	// 		Description:   group.Description,
	// 		AllowedDefinitions: make([]models.DefinitionRaw, len(group.AllowedDefinitions)),
	// 	}

	// 	for i, allowedModel := range group.AllowedModels {
	// 		groupExt.AllowedModels[i] = modelNameMap[allowedModel]
	// 	}

	// 	validGroups = append(validGroups, groupExt)

	// }

	// If all validations pass, return the groups
	return validGroups, nil

}

// Note: Reactive entity validation is performed at API time when entities are created/updated
// via API endpoints, not during initial configuration parsing.
