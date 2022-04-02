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

	idsQueryParam := c.Query("ids")
	// idsStrings := strings.Split(idsQueryParam, ",")

	// ids := make([]int, len(idsStrings))
	// for i := 0; i < len(idsQueryParam); i++ {

	// 	id, err := strconv.Atoi(idsStrings[i])
	// 	util.CheckErr(err)
	// 	ids[i] = id
	// }
	fmt.Printf("calling GetTSP start=%d ids=%s", start, idsQueryParam) //%d ids

	idsQueryParam = "55225524,33997995,33176384,240122791"
	tsp, err := routing.CalculateTSP(start, idsQueryParam) //ids
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
