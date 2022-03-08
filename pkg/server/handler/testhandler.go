package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	data "github.com/andidroid/testgo/pkg/server/data"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "data"
)

// Tests is a http.Handler
type TestHandler struct {
	ctx             context.Context
	log             *log.Logger
	mongoCollection *mongo.Collection
	redisClient     *redis.Client
}

// NewTestHandler creates a products handler with the given logger
func NewTestHandler(ctx context.Context, l *log.Logger, mongoDatabase *mongo.Database, redisClient *redis.Client) *TestHandler {
	mongoCollection := mongoDatabase.Collection("test")

	// cur, err := mongoCollection.Find(ctx, bson.M{})
	// defer cur.Close(ctx)
	// log.Printf("error: %s", err)
	// tests := make([]data.Test, 0)
	// for cur.Next(ctx) {
	// 	var test data.Test
	// 	cur.Decode(&test)
	// 	log.Printf("read %s", test)
	// 	tests = append(tests, test)
	// }

	return &TestHandler{ctx, l, mongoCollection, redisClient}
}

// swagger:operation GET /tests tests listTests
// Returns list of tests
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
// getTestHandler returns the products from the data store
func (handler *TestHandler) HandleGetAllTestsRequest(c *gin.Context) {
	handler.log.Println("Handle GET Tests")

	val, err := handler.redisClient.Get(handler.ctx, "tests").Result()
	if err == redis.Nil || err != nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.mongoCollection.Find(handler.ctx, bson.M{})
		if err != nil {
			log.Fatalf("Request to MongoDB failed: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)

		tests := make([]data.Test, 0)
		for cur.Next(handler.ctx) {
			var test data.Test
			cur.Decode(&test)
			tests = append(tests, test)
		}

		data, _ := json.Marshal(tests)
		s := handler.redisClient.Set(handler.ctx, "tests", string(data), time.Duration(1)*time.Minute)
		log.Printf("Set tests in redis cache %s", s)
		c.JSON(http.StatusOK, tests)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("Request to Redis")
		tests := make([]data.Test, 0)
		json.Unmarshal([]byte(val), &tests)
		c.JSON(http.StatusOK, tests)
	}
}

// swagger:operation GET /tests/{id} tests
// Get one test
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: test ID
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
func (handler *TestHandler) HandleGetTestByIdRequest(c *gin.Context) {
	handler.log.Println("Handle GET Test", c.Request.URL.String())
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.mongoCollection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var test data.Test
	err := cur.Decode(&test)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, test)
}

// swagger:operation POST /tests tests newRecipe
// Create a new test
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
func (handler *TestHandler) HandlePostRequest(c *gin.Context) {
	handler.log.Println("Handle POST Test")
	var test data.Test
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	test.ID = primitive.NewObjectID()
	test.RunAt = time.Now()
	_, err := handler.mongoCollection.InsertOne(handler.ctx, test)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new test"})
		return
	}

	log.Println("Remove data from Redis")
	handler.redisClient.Del(handler.ctx, "tests")

	c.JSON(http.StatusOK, test)
}

// swagger:operation PUT /tests/{id} tests updateRecipe
// Update an existing test
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the test
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
//     '404':
//         description: Invalid test ID
func (handler *TestHandler) HandlePutRequest(c *gin.Context) {
	id := c.Param("id")
	var test data.Test
	if err := c.ShouldBindJSON(&test); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.mongoCollection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", test.Name},
		{"description", test.Description},
		{"type", test.TestType},
		{"run", test.RunAt},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test has been updated"})
}

// swagger:operation DELETE /tests/{id} tests deleteRecipe
// Delete an existing test
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the test
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid test ID
func (handler *TestHandler) HandleDeleteRequest(c *gin.Context) {
	handler.log.Println("Handle DELETE Test", c.Request.URL.String())
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.mongoCollection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Test has been deleted"})
}
