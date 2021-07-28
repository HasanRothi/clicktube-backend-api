package db

import (
	"context"
	"fmt"
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
func Connect() {

	/*
	   Connect to my cluster
	*/
	DB_USER := loadDotEnvVariable("DB_USER")
	DB_PASSWORD := loadDotEnvVariable("DB_PASSWORD")
	DB_NAME := loadDotEnvVariable("DB_NAME")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + DB_USER + ":" + DB_PASSWORD + "@cluster0.gelqk.mongodb.net/" + DB_NAME + "?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	/*
	   List databases
	*/
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

}
