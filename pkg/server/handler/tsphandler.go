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
func GetTSP(c *gin.Context) {
	start, err := strconv.Atoi(c.Query("start"))
	util.CheckErr(err)

	fmt.Printf("calling GetTSP start=%d", start)

	tsp, err := routing.CalculateTSP(start)
	if err != nil {
		c.Status(500)
		return
	}

	c.JSON(http.StatusOK, tsp)
}

// func checkErr(err error) {
// 	if err != nil {
// 		fmt.Println("err connecting to  database", err)
// 	}
// }
