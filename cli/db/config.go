package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadDotEnvVariable(key string) string {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

var Database string
var DbClient *mongo.Client
var DbCtx context.Context

func Connect() {

	/*
	   Connect to cluster
	*/
	DB_USER := loadDotEnvVariable("DB_USER")
	DB_PASSWORD := loadDotEnvVariable("DB_PASSWORD")
	DB_NAME := loadDotEnvVariable("DB_NAME")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + DB_USER + ":" + DB_PASSWORD + "@cluster0.gelqk.mongodb.net/" + DB_NAME + "?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10000000*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	DbClient = client
	DbCtx = ctx
	// defer client.Disconnect(ctx)

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	Database = databases[0]

}
