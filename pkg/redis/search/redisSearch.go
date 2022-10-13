package search

import (
	"fmt"
	"log"
	"os"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
)

var redisearchClient *redisearch.Client

func init() {
	redisHost, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		redisHost = "localhost"
		os.Setenv("REDIS_HOST", "localhost")
	}
	fmt.Printf("REDIS_HOST: %s\n", redisHost)

	redisPort, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		redisPort = "6379"
		os.Setenv("REDIS_PORT", "6379")
	}
	fmt.Printf("REDIS_PORT: %s\n", redisPort)

	redisIndex, ok := os.LookupEnv("REDIS_INDEX")
	if !ok {
		redisIndex = "placesIndex"
		os.Setenv("REDIS_INDEX", "placesIndex")
	}
	fmt.Printf("REDIS_INDEX: %s\n", redisIndex)

	//redis://<user>:<password>@<host>:<port>/<db_number>
	redisURL := os.ExpandEnv("$REDIS_HOST:$REDIS_PORT")
	fmt.Println("Redis URL: ", redisURL)

	redisearchClient = redisearch.NewClient(redisURL, redisIndex)
}

func CreateClient() *redisearch.Client {
	return redisearchClient
}

func GetClient() *redisearch.Client {
	return redisearchClient
}

func TestRedisSearch() {

	fmt.Println("init redis search")
	redisearchClient = redisearch.NewClient("localhost:6379", "placesIndex")

	// Drop an existing index. If the index does not exist an error is returned
	//redisearchClient.Drop()

	redisearchClient.DropIndex(true)

	fmt.Println(redisearchClient)

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("id")).
		AddField(redisearch.NewSortableTextField("name", 5.0)).
		AddField(redisearch.NewNumericField("osm_id")).
		AddField(redisearch.NewTextField("description")).
		AddField(redisearch.NewGeoField("geom"))

	err := redisearchClient.CreateIndex(sc)
	util.CheckErr(err)

	features, err := routing.ReadAllPOIs()
	util.CheckErr(err)

	f := *features.Features
	for i := 0; i < len(f); i++ {
		feature := f[i]

		place := feature.Properties.(*routing.POI)

		doc := redisearch.NewDocument(string(place.ID), 1.0)
		doc.Set("name", place.Name).
			Set("geom", fmt.Sprintf("%f,%f", feature.Geom.Coordinates[0], feature.Geom.Coordinates[1])).
			//Set("description", place.Description).
			Set("id", place.ID).
			Set("osm_id", place.Osm_ID)
		//Set("date", time.Now().Unix())
		fmt.Println("insert ", doc, feature.ID)
		err := redisearchClient.IndexOptions(redisearch.DefaultIndexingOptions, doc)
		util.CheckErr(err)
	}

	// Index the document. The API accepts multiple documents at a time
	if err != nil {
		log.Fatal(err)
	}

	docs, total, err := redisearchClient.Search(redisearch.NewQuery("*").
		AddFilter(
			redisearch.Filter{
				Field: "geom",
				Options: redisearch.GeoFilterOptions{
					Lon:    12,
					Lat:    52,
					Radius: 100,
					Unit:   redisearch.KILOMETERS,
				},
			},
		).
		Limit(0, 10).
		SetReturnFields("name", "osm_id"))
	util.CheckErr(err)

	fmt.Println(docs)
	fmt.Println(total)
	fmt.Println(docs[0].Id, docs[0].Properties["name"], total, err)
}
