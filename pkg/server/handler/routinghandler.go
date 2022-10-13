package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andidroid/testgo/pkg/pgsql"
	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
	"github.com/blockloop/scan"
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
	util.CheckErr(err)
	target, err := strconv.Atoi(c.Query("target"))
	util.CheckErr(err)

	fmt.Printf("calling GetRouteAsList source=%d target=%d", source, target)

	conn, err := pgsql.GetConnection()
	util.CheckErr(err)

	sql := fmt.Sprintf("select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost FROM public.view_routing', %d , %d ,%t)", source, target, routing.DIRECTED_GRAPH)
	stmt, err := conn.Prepare(sql)
	util.CheckErr(err)
	//source=1529&target=2225
	//(source, target

	rows, err := stmt.Query() // (source, target)
	util.CheckErr(err)

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
	util.CheckErr(err)
	target, err := strconv.Atoi(c.Query("target"))
	util.CheckErr(err)
	format := c.Query("format")

	fmt.Printf("calling GetRouteAsGeometry source=%d target=%d format=%s", source, target, format)

	conn, err := pgsql.GetConnection()
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	var encoding string
	if format == "" {
		format = "ST_AsGeoJson"
		encoding = "application/json; charset=utf-8"
	} else if format == "geojson" {
		encoding = "application/json; charset=utf-8"
		format = "ST_AsGeoJson"
	} else if format == "wkt" {
		format = "ST_AsText"
		encoding = "application/text; charset=utf-8"
	}

	//ST_AsText
	//ST_Union
	sql := fmt.Sprintf("select %s(ST_LineMerge(ST_Collect(geom))) from public.roads where osm_id in (select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost FROM public.view_routing', %d , %d ,%t))", format, source, target, routing.DIRECTED_GRAPH)
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

	// fmt.Println("edge", edge)
	stmt.Close()

	c.Data(200, encoding, []byte(edge))
	//c.JSON(http.StatusOK, edge)
}

type PgrDijkstraRoute struct {
	Source int               `json:"source" db:"source"`
	Target int               `json:"target" db:"target"`
	Cost   float64           `json:"cost" db:"cost"`
	Steps  []PgrDijkstraStep `json:"steps" db:"steps"`
}

type PgrDijkstraStep struct {
	Seq     int     `json:"seq" db:"seq"`
	PathSeq int64   `json:"path_seq" db:"path_seq"`
	Node    int64   `json:"node" db:"node"`
	Edge    int64   `json:"edge" db:"edge"`
	Cost    float64 `json:"cost" db:"cost"`
	AggCost float64 `json:"agg_cost" db:"agg_cost"`
}

func GetRouteInformation(c *gin.Context) {
	source, err := strconv.Atoi(c.Query("source"))
	util.CheckErr(err)
	target, err := strconv.Atoi(c.Query("target"))
	util.CheckErr(err)

	fmt.Printf("calling GetRouteAsList source=%d target=%d", source, target)

	conn, err := pgsql.GetConnection()
	util.CheckErr(err)

	sql := fmt.Sprintf("select seq,path_seq,node,edge,cost,agg_cost FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost FROM public.view_routing', %d , %d ,%t)", source, target, routing.DIRECTED_GRAPH)
	// seq | path_seq | node | edge | cost | agg_cost
	// -----+----------+------+------+------+----------
	//    1 |        1 |    2 |    4 |    1 |        0
	//    2 |        2 |    5 |    8 |    1 |        1
	stmt, err := conn.Prepare(sql)
	util.CheckErr(err)
	//source=1529&target=2225
	//(source, target

	rows, err := stmt.Query() // (source, target)
	util.CheckErr(err)

	// rows, err := conn.Query("select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost, cost as reverse_cost FROM public.roads', $1 , $2 ,true)")
	// checkErr(err)

	// edges := make([]PgrDijkstraStep, 0, 10)
	// for rows.Next() {
	// 	var edge PgrDijkstraStep
	// 	err = rows.Scan(&edge)
	// 	if err != nil {
	// 		fmt.Println("err connecting to  database", err)
	// 	}
	// 	// checkErr(err)
	// 	fmt.Println("edge", edge)
	// 	// append(edges, edge)
	// 	edges = append(edges, edge)
	// }
	var steps []PgrDijkstraStep
	err = scan.Rows(&steps, rows)
	fmt.Println("steps", steps)
	rows.Close()
	stmt.Close()

	route := PgrDijkstraRoute{
		Source: source,
		Target: target,
		Steps:  steps,
	}

	c.JSON(http.StatusOK, route)
}
