package db

import (
	"fmt"

	"github.com/go-redis/redis"
)

// func loadDotEnvVariable(key string) string {

// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}
// 	return os.Getenv(key)
// }

var DbRedisClient *redis.Client

func RedisConnect() {

	DbRedisClient = redis.NewClient(&redis.Options{
		Addr:     loadDotEnvVariable("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
	})

	pong, err := DbRedisClient.Ping().Result()
	fmt.Println(pong, err)

}
