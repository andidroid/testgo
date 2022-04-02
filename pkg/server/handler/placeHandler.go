package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "data"
)

const distance int = 10000

type PlaceHandler struct {
	ctx             context.Context
	log             *log.Logger
	mongoCollection *mongo.Collection
	redisClient     *redis.Client
}

func NewPlaceHandler(ctx context.Context, l *log.Logger, mongoDatabase *mongo.Database, redisClient *redis.Client) *PlaceHandler {
	mongoCollection := mongoDatabase.Collection("places")
	return &PlaceHandler{ctx, l, mongoCollection, redisClient}
}

func (handler *PlaceHandler) HandleGetAllPlacesRequest(c *gin.Context) {
	log.Printf("HandleGetAllPlacesRequest")

	cur, err := handler.mongoCollection.Find(handler.ctx, bson.M{})
	if err != nil {
		log.Fatalf("Request to MongoDB failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	places := make([]routing.Place, 0)
	for cur.Next(handler.ctx) {
		var place routing.Place
		cur.Decode(&place)
		places = append(places, place)
	}
	c.JSON(http.StatusOK, places)
}

func (handler *PlaceHandler) HandleGetPlaceByIdRequest(c *gin.Context) {

	handler.log.Println("HandleGetPlaceByIdRequest", c.Request.URL.String())
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.mongoCollection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var place routing.Place
	err := cur.Decode(&place)
	if err != nil {
		log.Fatalf("Request to MongoDB failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	routing.SendMongoEvent(objectId.Hex(), "places", routing.MONGO_EVENT_READ)
	c.JSON(http.StatusOK, place)
}

func (handler *PlaceHandler) HandleGetPlaceByCoordinatesRequest(c *gin.Context) {

	handler.log.Println("HandleGetPlaceByCoordinatesRequest", c.Request.URL.String())

	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	util.CheckErr(err)
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	util.CheckErr(err)

	location := routing.Geometry{
		Type:        "Point",
		Coordinates: []float64{lon, lat},
	}
	var results []routing.Place
	filter := bson.D{
		{"location",
			bson.D{
				{"$near", bson.D{
					{"$geometry", location},
					{"$maxDistance", distance},
				}},
			}},
	}
	cur, err := handler.mongoCollection.Find(handler.ctx, filter)
	util.CheckErr(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for cur.Next(handler.ctx) {
		var p routing.Place
		err := cur.Decode(&p)
		util.CheckErr(err)
		if err != nil {
			fmt.Println("Could not decode Point")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		results = append(results, p)
	}

	// objectId, _ := primitive.ObjectIDFromHex(id)
	// cur := handler.mongoCollection.FindOne(handler.ctx, bson.M{
	// 	"_id": objectId,
	// })
	// var place routing.Place
	// err := cur.Decode(&place)
	// if err != nil {
	// 	log.Fatalf("Request to MongoDB failed: %s", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusOK, results)
}

func (handler *PlaceHandler) HandlePostPlaceRequest(c *gin.Context) {

	handler.log.Println("HandlePostPlaceRequest")
	var place routing.Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//place.ID = primitive.NewObjectID()
	res, err := handler.mongoCollection.InsertOne(handler.ctx, place)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new place"})
		return
	}
	//handler.redisClient.Se
	routing.SendMongoEvent(fmt.Sprintf("%x", res.InsertedID), "places", routing.MONGO_EVENT_CREATED)
	c.JSON(http.StatusOK, gin.H{"result": res})

}

func (handler *PlaceHandler) HandlePutPlaceRequest(c *gin.Context) {
	id := c.Param("id")
	var place routing.Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	res, err := handler.mongoCollection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", place.Name},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	routing.SendMongoEvent(fmt.Sprintf("%x", res.UpsertedID), "places", routing.MONGO_EVENT_UPDATED)
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (handler *PlaceHandler) HandleDeletePlaceRequest(c *gin.Context) {
	handler.log.Println("Handle DELETE Test", c.Request.URL.String())
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	res, err := handler.mongoCollection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	routing.SendMongoEvent(fmt.Sprintf("%x", objectId), "places", routing.MONGO_EVENT_DELETED)
	c.JSON(http.StatusOK, gin.H{"result": res})
}
