package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"fmt"
	"log"
	"github.com/joho/godotenv"
	"os"
)

var ClientMongo = ConnectDb()
var MondoDb = "echo_framework2"
//var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017/"+MondoDb)


var CategoryCollection = ClientMongo.Database(MondoDb).Collection("category")



func ConnectDb() *mongo.Client {

	errorVars := godotenv.Load()
	if errorVars != nil {
		panic("Error loading .env file")
	}
	var clientOptions = options.Client().ApplyURI("mongodb+srv://"+os.Getenv("MONGO_USER")+":"+os.Getenv("MONGO_PASSWORD")+"@clustertwittor.9efv4zw.mongodb.net/?retryWrites=true&w=majority&appName=ClusterTwittor")


	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
		return client
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func ConfirmConnection () int {
	err := ClientMongo.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}