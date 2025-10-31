// fetch.go
package persistence

import (
	"databus/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Assuming these are the collection names you are using:
const (
	definitionsCollection     = "Definitions"
	groupsCollection          = "Groups"
	reactiveEntitiesCollection = "ReactiveEntities"
)

// GetAllDefinitions retrieves all models from the MongoDB collection "Models".
func GetAllDefinitions() ([]models.DefinitionRaw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(definitionsCollection)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.DefinitionRaw
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetDefinitionByID retrieves a single model by its ID from the MongoDB collection "Definitions".
func GetDefinitionByID(id primitive.ObjectID) (*models.DefinitionRaw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(definitionsCollection)
	var result models.DefinitionRaw
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDefinitionByName retrieves a single model by its name from the MongoDB collection "Definitions".
func GetDefinitionByName(name string) (*models.DefinitionRaw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(definitionsCollection)
	var result models.DefinitionRaw
	err := collection.FindOne(ctx, bson.M{"Name": name}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllGroups retrieves all groups from the MongoDB collection "Groups".
func GetAllGroups() ([]models.GroupRaw, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(groupsCollection)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.GroupRaw
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetGroupByID retrieves a group by its ID from the MongoDB collection "Groups".
func GetGroupByID(id primitive.ObjectID) (*models.GroupRaw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(groupsCollection)
	var result models.GroupRaw
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetGroupByName retrieves a group by its name from the MongoDB collection "Groups".
func GetGroupByName(name string) (*models.GroupRaw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(groupsCollection)
	var result models.GroupRaw
	err := collection.FindOne(ctx, bson.M{"Name": name}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAllReactiveEntities retrieves all reactive entities from the MongoDB collection "ReactiveEntities".
func GetAllReactiveEntities() ([]models.ReactiveEntityRaw, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(reactiveEntitiesCollection)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ReactiveEntityRaw
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// GetReactiveEntityByID retrieves a reactive entity by its ID from the MongoDB collection "ReactiveEntities".
func GetReactiveEntityByID(id primitive.ObjectID) (*models.ReactiveEntityRaw, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(reactiveEntitiesCollection)
	var result models.ReactiveEntityRaw
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetReactiveEntityByHex(hex uint16) (*models.ReactiveEntityRaw, error) {
	// retrieves a reactive entity by its hex ID from the MongoDB collection "ReactiveEntities".

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(reactiveEntitiesCollection)
	var result models.ReactiveEntityRaw
	err := collection.FindOne(ctx, bson.M{"EntityHex": hex}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetReactiveEntitiesByGroup(groupsParam []string) ([]models.ReactiveEntityRaw, error) {
	// GetReactiveEntitiesByGroup retrieves all reactive entities that belong to the specified group(s).
	// groupsParam is a list of group names.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Retrieve group IDs corresponding to the provided group names
	groupCollection := MongoClient.Database("databus").Collection("Groups")
	groupCursor, err := groupCollection.Find(ctx, bson.M{"Name": bson.M{"$in": groupsParam}})
	if err != nil {
		return nil, err
	}
	defer groupCursor.Close(ctx)

	var groups []models.GroupRaw
	if err = groupCursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	groupIDs := make([]primitive.ObjectID, len(groups))
	for i, group := range groups {
		groupIDs[i] = group.ID
	}

	// Step 2: Query the reactive entities collection to find entities that belong to any of the retrieved group IDs
	reactiveEntityCollection := MongoClient.Database("databus").Collection("ReactiveEntities")
	reactiveEntityCursor, err := reactiveEntityCollection.Find(ctx, bson.M{"Groups": bson.M{"$all": groupIDs}})
	if err != nil {
		return nil, err
	}
	defer reactiveEntityCursor.Close(ctx)

	var results []models.ReactiveEntityRaw
	if err = reactiveEntityCursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// DeleteReactiveEntityByHex deletes a reactive entity by its hex ID
func DeleteReactiveEntityByHex(hex uint16) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := MongoClient.Database("databus").Collection(reactiveEntitiesCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"EntityHex": hex})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
