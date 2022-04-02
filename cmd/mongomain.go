package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main_m() {

	fmt.Println("init mongo database")

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

	mongodb := os.ExpandEnv("mongodb://$MONGODB_HOST:$MONGODB_PORT/")
	fmt.Println("MongoDB URL: ", mongodb)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	opts := options.Client()
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	mongoDatabase := mongoClient.Database("test")
	placesCollection := mongoDatabase.Collection("places")

	// Create a schema
	// sc := redisearch.NewSchema(redisearch.DefaultOptions).
	// 	AddField(redisearch.NewTextField("id")).
	// 	AddField(redisearch.NewSortableTextField("name", 5.0)).
	// 	AddField(redisearch.NewNumericField("osm_id")).
	// 	AddField(redisearch.NewTextField("description")).
	// 	AddField(redisearch.NewGeoField("geom"))

	features, err := routing.ReadAllPOIs()
	util.CheckErr(err)

	f := *features.Features
	for i := 0; i < len(f); i++ {
		feature := f[i]

		place := feature.Properties.(*routing.POI)

		result, err := placesCollection.InsertOne(ctx, bson.D{
			{Key: "id", Value: place.ID},
			{Key: "name", Value: place.Name},
			{Key: "geom", Value: feature.Geom},
			{Key: "osm_id", Value: place.Osm_ID},
		})
		util.CheckErr(err)
		fmt.Println(result)

	}

	// Index the document. The API accepts multiple documents at a time
	if err != nil {
		log.Fatal(err)
	}

}
