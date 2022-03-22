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
	// "data"
)

func HandleGetAllTrucksRequest(c *gin.Context) {
	conn, err := pgsql.GetConnection()
	util.CheckErr(err)

	sql := fmt.Sprintf("SELECT * FROM trucks")
	stmt, err := conn.Prepare(sql)
	util.CheckErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	util.CheckErr(err)
	defer rows.Close()

	var trucks []routing.Truck
	err = scan.Rows(&trucks, rows)

	c.JSON(http.StatusOK, trucks)
}

func HandleGetTruckRequest(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	util.CheckErr(err)

	conn, err := pgsql.GetConnection()
	util.CheckErr(err)

	sql := fmt.Sprintf("SELECT * FROM trucks WHERE id = %d", id)
	stmt, err := conn.Prepare(sql)
	util.CheckErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	util.CheckErr(err)
	defer rows.Close()
	var truck routing.Truck
	err = scan.Row(&truck, rows)

	c.JSON(http.StatusOK, truck)
}

func HandlePostTruckRequest(c *gin.Context) {

	var truck routing.Truck
	c.Bind(&truck)

	conn, err := pgsql.GetConnection()
	util.CheckErr(err)

	// ID         int64      `json:"id" bson:"_id"`
	// Name       string     `json:"name" bson:"name"`
	// Base       int64      `json:"base" bson:"base"`

	sql := "INSERT INTO trucks (name,base) VALUES ($2,$3) RETURNING id"
	stmt, err := conn.Prepare(sql)
	// var id int64
	res, err := stmt.Exec(sql, truck.Name, truck.Base)
	defer stmt.Close()
	util.CheckErr(err)
	id, err := res.LastInsertId()
	util.CheckErr(err)

	c.JSON(http.StatusOK, id)

}

func HandlePutTruckRequest(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	util.CheckErr(err)

	var truck routing.Truck
	c.Bind(&truck)

	conn, err := pgsql.GetConnection()
	util.CheckErr(err)

	// ID         int64      `json:"id" bson:"_id"`
	// Name       string     `json:"name" bson:"name"`
	// Base       int64      `json:"base" bson:"base"`

	// sql := fmt.Sprintf("UPDATE truck SET name=%s, base=%d  WHERE id = %d", id)
	sql := "UPDATE trucks SET name=$2, base=$3  WHERE id = $1"
	stmt, err := conn.Prepare(sql)
	res, err := stmt.Exec(sql, id, truck.Name, truck.Base)
	util.CheckErr(err)
	defer stmt.Close()
	c.JSON(http.StatusOK, res)
}

func HandleDeleteTruckRequest(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	util.CheckErr(err)

	conn, err := pgsql.GetConnection()
	util.CheckErr(err)

	sql := fmt.Sprintf("DELETE FROM trucks WHERE id = $1")
	stmt, err := conn.Prepare(sql)
	util.CheckErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(id)
	util.CheckErr(err)
	rowsAffected, err := res.RowsAffected()

	c.JSON(http.StatusOK, rowsAffected)

}
