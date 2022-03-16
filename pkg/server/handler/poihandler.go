package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
	"github.com/gin-gonic/gin"
)

// type RoutingHandler struct {
// }

// // NewTestHandler creates a products handler with the given logger
// func NewRoutingHandler() *RoutingHandler {
// 	return &RoutingHandler{}
// }

//localhost:8080/routing/list?source=1529&target=2225
func GetPOIs(c *gin.Context) {
	fmt.Printf("calling GetPOIs")

	pois, err := routing.ReadAllPOIs()
	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, &pois)
	}

}

func GetNearestNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("osm_id"))
	util.CheckErr(err)

	node, err := routing.FindNearestNodeForPOI(id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusOK, &node)
	}

}

// func checkErr(err error) {
// 	if err != nil {
// 		fmt.Println("err connecting to  database", err)
// 	}
// }
