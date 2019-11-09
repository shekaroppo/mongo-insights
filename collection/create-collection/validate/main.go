package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:single@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("noindex")

	result := database.RunCommand(context.Background(), bson.M{
		"create": "validation",
		"validator": bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": bson.A{"name", "age"},
				"properties": bson.M{
					"name": bson.M{
						"bsonType":    "string",
						"description": "must be a string and is required",
					},
					"age": bson.M{
						"bsonType":    "int",
						"description": "must be a int and is required",
					},
				},
			},
		},
	})

	if result.Err() != nil {
		raw, _ := result.DecodeBytes()
		log.Println(raw.String())
		return
	}

	raw, _ := result.DecodeBytes()
	log.Println(raw.String())
}
