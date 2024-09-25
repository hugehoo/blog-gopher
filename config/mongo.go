package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Client *mongo.Client

func ConnectMongoDB(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("🔥Mongo Connect Fail", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("🔥Mongo Connect Fail", err)
	}

	log.Println("MongoDB에 성공적으로 연결되었습니다!")
	Client = client
	return client
}

func GetCollection(database, collection string) *mongo.Collection {
	return Client.Database(database).Collection(collection)
}
