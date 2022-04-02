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

type OrderHandler struct {
	ctx             context.Context
	log             *log.Logger
	mongoCollection *mongo.Collection
	redisClient     *redis.Client
}

func NewOrderHandler(ctx context.Context, l *log.Logger, mongoDatabase *mongo.Database, redisClient *redis.Client) *OrderHandler {
	mongoCollection := mongoDatabase.Collection("orders")
	return &OrderHandler{ctx, l, mongoCollection, redisClient}
}

func (handler *OrderHandler) HandleGetAllOrdersRequest(c *gin.Context) {
	log.Printf("HandleGetAllOrdersRequest")

	cur, err := handler.mongoCollection.Find(handler.ctx, bson.M{})
	if err != nil {
		log.Fatalf("Request to MongoDB failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	orders := make([]routing.Order, 0)
	for cur.Next(handler.ctx) {
		var order routing.Order
		cur.Decode(&order)
		orders = append(orders, order)
	}
	c.JSON(http.StatusOK, orders)
}

func (handler *OrderHandler) HandleGetOrderByIdRequest(c *gin.Context) {

	handler.log.Println("HandleGetOrderByIdRequest", c.Request.URL.String())
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.mongoCollection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var order routing.Order
	err := cur.Decode(&order)
	if err != nil {
		log.Fatalf("Request to MongoDB failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (handler *OrderHandler) HandlePostOrderRequest(c *gin.Context) {

	handler.log.Println("HandlePostOrderRequest")
	var order routing.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = primitive.NewObjectID()
	res, err := handler.mongoCollection.InsertOne(handler.ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new order"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})

}

func (handler *OrderHandler) HandlePutOrderRequest(c *gin.Context) {
	id := c.Param("id")
	var order routing.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	res, err := handler.mongoCollection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", order.Name},
		{"starttime", order.StartTime},
		{"endtime", order.EndTime},
		{"state", order.OrderState},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (handler *OrderHandler) HandleDeleteOrderRequest(c *gin.Context) {
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
