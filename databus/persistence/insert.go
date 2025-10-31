// insert.go
package persistence

import (
	"databus/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertDefinitions(client *mongo.Client, definitions []models.DefinitionRaw) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	definitionCollection := MongoClient.Database("databus").Collection("Definitions")

	// Drop the collection (fresh start from incoming file)
	err := definitionCollection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping dataObjectCollection:", err)
	}

	definitionInterface := make([]interface{}, len(definitions))
	for i, def := range definitions {
		definitionInterface[i] = def
	}
	// Insert definitions into the collection
	result, err := definitionCollection.InsertMany(ctx, definitionInterface)
	if err != nil {
		return nil, fmt.Errorf("error inserting definitions: %v", err)
	}

	// Iterate through the inserted IDs and populate ModelIdMap
	for i, id := range result.InsertedIDs {
		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("unexpected ID type for definition at index %d", i)
		}

		definitions[i].ID = objectID
	}

	return result, nil
}

func InsertGroups(client *mongo.Client, groups []models.GroupRaw) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	groupCollection := MongoClient.Database("databus").Collection("Groups")

	// Drop the collection (fresh start from incoming file)
	err := groupCollection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping dataObjectCollection:", err)
	}

	groupInterfaces := make([]interface{}, len(groups))
	for i, group := range groups {
		groupInterfaces[i] = group
	}
	result, err := groupCollection.InsertMany(ctx, groupInterfaces)

	if err != nil {
		log.Fatal("Error inserting sample data:", err)
	}

	// Iterate through the inserted IDs and populate ModelIdMap
	for i, id := range result.InsertedIDs {
		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("unexpected ID type for definition at index %d", i)
		}

		groups[i].ID = objectID
	}

	return result, nil
}

func InsertReactiveEntities(client *mongo.Client, reactiveEntities []models.ReactiveEntityRaw) (*mongo.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reactiveEntityCollection := MongoClient.Database("databus").Collection("ReactiveEntities")

	// Drop the collection (fresh start from incoming file)
	err := reactiveEntityCollection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping reactiveEntityCollection:", err)
	}

	reactiveEntityInterfaces := make([]interface{}, len(reactiveEntities))
	for i, reactiveEntity := range reactiveEntities {
		reactiveEntityInterfaces[i] = reactiveEntity
	}
	result, err := reactiveEntityCollection.InsertMany(ctx, reactiveEntityInterfaces)

	if err != nil {
		log.Fatal("Error inserting reactive entities:", err)
	}

	// Iterate through the inserted IDs and populate IDs
	for i, id := range result.InsertedIDs {
		objectID, ok := id.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("unexpected ID type for reactive entity at index %d", i)
		}

		reactiveEntities[i].ID = objectID
	}
	return result, nil
}

// InsertReactiveEntity inserts a single reactive entity into the database (used by API)
func InsertReactiveEntity(reactiveEntity *models.ReactiveEntityRaw) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reactiveEntityCollection := MongoClient.Database("databus").Collection("ReactiveEntities")

	result, err := reactiveEntityCollection.InsertOne(ctx, reactiveEntity)
	if err != nil {
		return nil, fmt.Errorf("error inserting reactive entity: %v", err)
	}

	// Update the ID field with the inserted ID
	if insertedID, ok := result.InsertedID.(primitive.ObjectID); ok {
		reactiveEntity.ID = insertedID
	}

	return result, nil
}

func GetGroupIDMap(ctx context.Context, client *mongo.Client) (map[string]primitive.ObjectID, error) {
	groupCollection := client.Database("databus").Collection("groups")

	cursor, err := groupCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	groupNameToID := make(map[string]primitive.ObjectID)
	for cursor.Next(ctx) {
		var group models.GroupRaw
		if err := cursor.Decode(&group); err != nil {
			return nil, err
		}
		groupNameToID[group.Name] = group.ID
	}

	return groupNameToID, nil
}
