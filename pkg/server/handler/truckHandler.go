package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/andidroid/testgo/pkg/routing"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "data"
)

type TruckHandler struct {
	ctx             context.Context
	log             *log.Logger
	mongoCollection *mongo.Collection
	redisClient     *redis.Client
}

func NewTruckHandler(ctx context.Context, l *log.Logger, mongoDatabase *mongo.Database, redisClient *redis.Client) *TruckHandler {
	mongoCollection := mongoDatabase.Collection("trucks")
	return &TruckHandler{ctx, l, mongoCollection, redisClient}
}

func (handler *TruckHandler) HandleGetAllTrucksRequest(c *gin.Context) {
	log.Printf("HandleGetAllTrucksRequest")

	cur, err := handler.mongoCollection.Find(handler.ctx, bson.M{})
	if err != nil {
		log.Fatalf("Request to MongoDB failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	trucks := make([]routing.Truck, 0)
	for cur.Next(handler.ctx) {
		var truck routing.Truck
		cur.Decode(&truck)
		trucks = append(trucks, truck)
	}
	c.JSON(http.StatusOK, trucks)
}

func (handler *TruckHandler) HandleGetTruckByIdRequest(c *gin.Context) {

	handler.log.Println("HandleGetTruckByIdRequest", c.Request.URL.String())
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.mongoCollection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var truck routing.Truck
	err := cur.Decode(&truck)
	if err != nil {
		log.Fatalf("Request to MongoDB failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, truck)
}

func (handler *TruckHandler) HandlePostTruckRequest(c *gin.Context) {

	handler.log.Println("HandlePostTruckRequest")
	var truck routing.Truck
	if err := c.ShouldBindJSON(&truck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	truck.ID = primitive.NewObjectID()
	res, err := handler.mongoCollection.InsertOne(handler.ctx, truck)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new truck"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})

}

func (handler *TruckHandler) HandlePutTruckRequest(c *gin.Context) {
	id := c.Param("id")
	var truck routing.Truck
	if err := c.ShouldBindJSON(&truck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	res, err := handler.mongoCollection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", truck.Name},
		{"base", truck.Base},
		{"action", truck.ActionType},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (handler *TruckHandler) HandleDeleteTruckRequest(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"result": res})
}
