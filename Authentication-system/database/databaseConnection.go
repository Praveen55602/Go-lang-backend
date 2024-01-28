package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

// trying to connect to the mongoDB
func DBinstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//In go all the api requests and responses are handled by an object called context. Unlike express which has separate object of request and response object gin has a object called context which handles both.
	//here we are creating a new context object with a timeout of 10 seconds that will be send to the connect command to store all the details regarding the call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//cancel function cleans up the resources associated with the context and will be called once the this functions exits or if the context takes more time then the timeout timer, until then it will be defered
	defer cancel()

	MongoDb := os.Getenv("MONGODB_URL")
	clientOptions := options.Client().ApplyURI(MongoDb)
	//passing context object where all the result of api call will be put
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	//checking the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	fmt.Println("successfully connected to mongoDB")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client mongo.Client, collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DB_NAME")
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}
