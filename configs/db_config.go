package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	uri := os.Getenv("DB_URI")
	fmt.Println("DB_URI ", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	fmt.Println("Success connected db")
	Client = client
	return client, nil
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DB_NAME")
	fmt.Println("DB_NAME: ", databaseName)

	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}
