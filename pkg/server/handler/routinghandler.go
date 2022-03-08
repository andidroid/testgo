package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andidroid/testgo/pkg/pgsql"
	"github.com/gin-gonic/gin"
)

// type RoutingHandler struct {
// }

// // NewTestHandler creates a products handler with the given logger
// func NewRoutingHandler() *RoutingHandler {
// 	return &RoutingHandler{}
// }

//localhost:8080/routing/list?source=1529&target=2225
func GetRouteAsList(c *gin.Context) {
	source, err := strconv.Atoi(c.Query("source"))
	checkErr(err)
	target, err := strconv.Atoi(c.Query("target"))
	checkErr(err)

	fmt.Printf("calling GetRouteAsList source=%d target=%d", source, target)

	conn, err := pgsql.InitDB()
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	sql := fmt.Sprintf("select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost FROM public.view_routing', %d , %d ,true)", source, target)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}
	//source=1529&target=2225
	//(source, target

	rows, err := stmt.Query() // (source, target)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	// rows, err := conn.Query("select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost, cost as reverse_cost FROM public.roads', $1 , $2 ,true)")
	// checkErr(err)

	edges := make([]int64, 0, 10)
	for rows.Next() {
		var edge int64
		err = rows.Scan(&edge)
		if err != nil {
			fmt.Println("err connecting to  database", err)
		}
		// checkErr(err)
		fmt.Println("edge", edge)
		// append(edges, edge)
		edges = append(edges, edge)
	}
	fmt.Println("edges", edges)
	rows.Close()
	stmt.Close()

	c.JSON(http.StatusOK, edges)
}

//localhost:8080/routing/geometry?source=1529&target=2225
func GetRouteAsGeometry(c *gin.Context) {
	source, err := strconv.Atoi(c.Query("source"))
	checkErr(err)
	target, err := strconv.Atoi(c.Query("target"))
	checkErr(err)

	fmt.Printf("calling GetRouteAsGeometry source=%d target=%d", source, target)

	conn, err := pgsql.InitDB()
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	sql := fmt.Sprintf("select ST_AsText(ST_LineMerge(ST_Union(geom))) from public.roads where osm_id in (select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost FROM public.view_routing', %d , %d ,true))", source, target)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}
	//source=1529&target=2225
	//(source, target

	var edge string
	err = stmt.QueryRow().Scan(&edge) // (source, target)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	// rows, err := conn.Query("select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost, cost as reverse_cost FROM public.roads', $1 , $2 ,true)")
	// checkErr(err)

	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	fmt.Println("edge", edge)
	stmt.Close()

	c.JSON(http.StatusOK, edge)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}
}
