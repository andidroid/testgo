package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/andidroid/testgo/pkg/redis"
	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
	"github.com/gin-gonic/gin"
	// "data"
)

type SearchPlaceHandler struct {
	ctx              context.Context
	log              *log.Logger
	redisearchClient *redisearch.Client
}

func NewSearchPlaceHandler(ctx context.Context, l *log.Logger, redisearchClient *redisearch.Client) *SearchPlaceHandler {
	initSubscriber(ctx, redisearchClient)
	return &SearchPlaceHandler{ctx, l, redisearchClient}
}

func (handler *SearchPlaceHandler) HandleGetSearchRequest(c *gin.Context) {

	handler.log.Println("HandleGetSearchRequest", c.Request.URL.String())
	query := c.Query("query")

	docs, total, err := handler.redisearchClient.Search(redisearch.NewQuery(query).
		Limit(0, 10).
		SetReturnFields("name", "osm_id"))
	util.CheckErr(err)

	fmt.Println(docs)
	fmt.Println(total)
	fmt.Println(docs[0].Id, docs[0].Properties["name"], total, err)

	if err != nil {
		log.Fatalf("Request to Redis Search failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

func initSubscriber(ctx context.Context, redisearchClient *redisearch.Client) {
	go func() {
		sub := redis.GetClient().Subscribe(ctx, "mongoevent")
		pe := routing.MongoEvent{}
		for {
			msg, err := sub.ReceiveMessage(ctx)
			util.CheckErr(err)
			fmt.Println(msg)

			if err := json.Unmarshal([]byte(msg.Payload), &pe); err != nil {
				util.CheckErr(err)
			} else {
				// fmt.Println("received mongo event: ", pe)
				switch pe.Type {
				case routing.MONGO_EVENT_READ:
					fmt.Println("read mongo event: ", pe)
				case routing.MONGO_EVENT_CREATED:
					fmt.Println("created mongo event: ", pe)

					url := fmt.Sprintf("http://localhost/fleet/place/%s", pe.ID)
					resp, err := routing.GetClient().Get(url)
					util.CheckErr(err)
					fmt.Println("server GET place: ", pe.ID, resp, err)

					defer resp.Body.Close()
					nodeBytes, err := ioutil.ReadAll(resp.Body)
					util.CheckErr(err)
					//nodeString := string(nodeBytes)
					place := routing.Place{}
					err = json.Unmarshal(nodeBytes, &place)
					util.CheckErr(err)

					doc := redisearch.NewDocument(place.ID.String(), 1.0)
					doc.Set("name", place.Name).
						Set("geom", fmt.Sprintf("%f,%f", place.Geometry.Coordinates[0], place.Geometry.Coordinates[1])).
						//Set("description", place.Description).
						Set("id", place.ID).
						Set("osm_id", place.OsmId)
					//Set("date", time.Now().Unix())
					err = redisearchClient.IndexOptions(redisearch.DefaultIndexingOptions, doc)
					util.CheckErr(err)

				default:
					fmt.Println("error receiving mongo event: ", pe)
				}
			}

		}
	}()
}
