package db

import (
	"context"
	"fmt"
	"log"

	"backend/env"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.Background()
var RDB *redis.Client
var Client *mongo.Client

// collections
var Accounts *mongo.Collection
var TicketTypes *mongo.Collection
var Tickets *mongo.Collection

func InitDB() {
	var err error

	Client, err = mongo.Connect(
		Ctx,
		options.Client().ApplyURI(env.MongoURI),
	)

	if err != nil {
		log.Fatal(err)
	}

	// loading collections
	Accounts = GetCollection("accounts", Client)
	TicketTypes = GetCollection("tickettypes", Client)
	Tickets = GetCollection("tickets", Client)

	fmt.Println("connected to mongodb")
}

func GetCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	return client.Database("dev").Collection(collectionName)
}

func InitCache() {
	RDB = redis.NewClient(env.RedisOptions)

	pong, _ := RDB.Ping(context.Background()).Result()
	if pong == "PONG" {
		fmt.Println("connected to redis")
	} else {
		fmt.Println("not connected to redis")
	}
}

func Set(key string, value string) error {
	err := RDB.Set(Ctx, key, value, 0).Err()

	return err
}

func Get(key string) (string, error) {
	val, err := RDB.Get(Ctx, key).Result()

	return val, err
}

func Del(key string) error {
	_, err := RDB.Del(Ctx, key).Result()

	return err
}
