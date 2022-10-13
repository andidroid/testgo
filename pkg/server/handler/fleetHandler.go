package handler

import (
	"github.com/andidroid/testgo/pkg/routing"
	"github.com/gin-gonic/gin"
	// "data"
)

func HandleGetFleetRequest(c *gin.Context) {
	//var fleet routing.Fleet

}

func HandlePostFleetRequest(c *gin.Context) {

}

func HandlePostOrderRequest(c *gin.Context) {

}

type StartOrderRequest struct {
	Order string `form:"order" json:"order" bson:"order" binding:"required"`
	Truck string `form:"truck" json:"truck" bson:"truck" binding:"required"`
	Start int    `form:"start" json:"start" bson:"start" binding:"required"`
}

func HandlePostStartOrderRequest(c *gin.Context) {
	var StartOrderRequest StartOrderRequest
	c.Bind(&StartOrderRequest)
	routing.GetFleetInstance().StartOrder(StartOrderRequest.Order, StartOrderRequest.Truck, StartOrderRequest.Start)
}
