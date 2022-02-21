package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Main() {

	mongodbHost, ok := os.LookupEnv("MONGODB_HOST")
	if !ok {
		mongodbHost = "localhost"
		os.Setenv("MONGODB_HOST", "localhost")
	}
	fmt.Printf("MONGODB_HOST: %s\n", mongodbHost)

	mongodbPort, ok := os.LookupEnv("MONGODB_PORT")
	if !ok {
		mongodbPort = "27017"
		os.Setenv("MONGODB_PORT", "27017")
	}
	fmt.Printf("MONGODB_PORT: %s\n", mongodbPort)

	mongodb := os.ExpandEnv("mongodb://$MONGODB_HOST:$MONGODB_PORT")
	fmt.Println("MongoDB URL: ", mongodb)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb))
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("test")
	res, err := collection.InsertOne(ctx, bson.D{{"name", "test"}, {"value", 3.14159}})
	id := res.InsertedID
	if id != nil {
		fmt.Println(id)
	}

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		fmt.Println(result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}
