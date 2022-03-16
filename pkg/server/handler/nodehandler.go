package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andidroid/testgo/pkg/pgsql"
	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
	"github.com/gin-gonic/gin"
)

type NodeType int

const (
	SOURCE NodeType = iota + 1
	TARGET
	UNKNOWN
)

func GetNodeSearchQuery(c *gin.Context) {
	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	util.CheckErr(err)
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	util.CheckErr(err)
	query := c.Query("query")

	fmt.Printf("calling NodeSearchQuery lng=%f lat=%f query=%s", lon, lat, query)

	node := findNode(lon, lat, UNKNOWN)

	c.JSON(http.StatusOK, node)
}

func GetNodeSearchSource(c *gin.Context) {
	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	util.CheckErr(err)
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	util.CheckErr(err)
	query := c.Query("query")

	fmt.Printf("calling NodeSearchQuery lng=%f lat=%f query=%s", lon, lat, query)

	node := findNode(lon, lat, SOURCE)

	c.JSON(http.StatusOK, node)
}

func GetNodeSearchTarget(c *gin.Context) {
	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	util.CheckErr(err)
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	util.CheckErr(err)
	query := c.Query("query")

	fmt.Printf("calling NodeSearchQuery lng=%f lat=%f query=%s", lon, lat, query)

	node := findNode(lon, lat, TARGET)

	c.JSON(http.StatusOK, node)
}

func GetNodeById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	util.CheckErr(err)

	fmt.Printf("calling NodeSearchQuery id=%d ", id)

	node := findNodeById(id)

	c.JSON(http.StatusOK, node)
}

func findNode(lon float64, lat float64, nodeType NodeType) routing.Node {

	conn, err := pgsql.GetConnection()
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	filter := ""
	switch nodeType {
	case SOURCE:
		filter = "where eout > 0"
		break
	case TARGET:
		filter = "where ein > 0"
		break
	default:
		filter = ""
		break
	}

	// find nearest node
	sql := fmt.Sprintf("select id, ST_X(the_geom),ST_Y(the_geom) from public.roads_vertices_pgr "+filter+" order by the_geom <-> ST_GeomFromText( 'POINT(%f %f)' , 4326) limit 1", lon, lat)
	fmt.Println(sql)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	var id int64
	err = stmt.QueryRow().Scan(&id, &lon, &lat)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	fmt.Printf("node id:%d, lon:%f, lat%f", id, lon, lat)
	fmt.Println("node id:", id)
	stmt.Close()

	node := routing.Node{
		ID: id, Lon: lon, Lat: lat,
	}
	return node
}

func findNodeById(id int64) routing.Node {

	conn, err := pgsql.GetConnection()
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	// find nearest node
	sql := fmt.Sprintf("select id, ST_X(the_geom),ST_Y(the_geom) from public.roads_vertices_pgr where id=%d", id)
	fmt.Println(sql)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	var lon float64
	var lat float64
	err = stmt.QueryRow().Scan(&id, &lon, &lat)
	if err != nil {
		fmt.Println("err connecting to  database", err)
	}

	fmt.Printf("node id:%d, lon:%f, lat%f", id, lon, lat)
	stmt.Close()

	node := routing.Node{
		ID: id, Lon: lon, Lat: lat,
	}
	return node
}
