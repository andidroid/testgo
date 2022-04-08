package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	// "data"
)

// Tests is a http.Handler
type HealthHandler struct {
	ctx           context.Context
	log           *log.Logger
	mongoDatabase *mongo.Database
	redisClient   *redis.Client
}

type Status struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Substatus []SubStatus
}

type SubStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NewHealthHandler creates a products handler with the given logger
func NewHealthHandler(ctx context.Context, l *log.Logger, mongoDatabase *mongo.Database, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{ctx, l, mongoDatabase, redisClient}
}

// swagger:operation GET /tests tests listTests
// Returns list of tests
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
// getHealthHandler returns the products from the data store
func (handler *HealthHandler) HandleGetRequest(c *gin.Context) {
	handler.log.Println("Handle GET health")

	foundError := false
	subStatusList := make([]SubStatus, 0)

	redisPing := handler.redisClient.Ping(handler.ctx)
	redisPingResult, err := redisPing.Result()
	err = redisPing.Err()

	// err != redis.Nil ||
	if err != nil {
		foundError = true
		// subStatus := SubStatus{"DOWN", redisPingResult}
		var subStatus SubStatus
		subStatus.Name = "Redis Client"
		subStatus.Status = "DOWN"
		subStatus.Message = redisPingResult
		subStatusList = append(subStatusList, subStatus)
	} else {
		//subStatus := SubStatus{"UP", redisPingResult}
		var subStatus SubStatus
		subStatus.Name = "Redis Client"
		subStatus.Status = "UP"
		subStatus.Message = redisPingResult
		subStatusList = append(subStatusList, subStatus)
	}

	err = handler.mongoDatabase.Client().Ping(handler.ctx, nil)
	if err != nil {
		foundError = true
		subStatus := SubStatus{"MongoDB", "DOWN", ""}
		subStatusList = append(subStatusList, subStatus)
	} else {
		subStatus := SubStatus{"MongoDB", "UP", ""}
		subStatusList = append(subStatusList, subStatus)
	}

	if foundError {
		// status := Status{"DOWN", "", subStatusList}
		var status Status
		status.Status = "DOWN"
		status.Message = ""
		status.Substatus = subStatusList
		fmt.Printf("health %s", status)
		c.JSON(http.StatusInternalServerError, status)
	} else {
		// status := Status{"UP", "", subStatusList}
		var status Status
		status.Status = "UP"
		status.Message = ""
		status.Substatus = subStatusList
		fmt.Printf("health %s", status)
		c.JSON(http.StatusInternalServerError, status)
	}
}
