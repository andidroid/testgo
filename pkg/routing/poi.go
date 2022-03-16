package routing

import (
	"encoding/json"
	"fmt"

	"github.com/andidroid/testgo/pkg/pgsql"
	"github.com/andidroid/testgo/pkg/util"
	"github.com/blockloop/scan"
)

type POI struct {
	ID       int64  `json:"id" bson:"_id" db:"fid"`
	Osm_ID   string `json:"osm_id" bson:"osm_id" db:"osm_id"`
	Name     string `json:"name" bson:"name"`
	GeomJSON string `db:"geom" json:"-"`
}

//geom,omitempty

//{"type":"Point","coordinates":[12.4045328,51.7979734]}
type GeoJSON struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type FeatureCollection struct {
	Type    string     `json:"type" bson:"type"`
	Feature *[]Feature `json:"features" bson:"features"`
	Name    string     `json:"name" bson:"name"`
}

type Feature struct {
	Type string  `json:"type" bson:"type"`
	Geom GeoJSON `json:"geometry" bson:"geom" db:"-"`
	POI  *POI    `json:"properties" bson:"properties"`
}

func ReadAllPOIs() (*FeatureCollection, error) {
	//func ReadAllPOIs() (*[]POI, error) {

	conn, err := pgsql.GetConnection()
	if err != nil {
		fmt.Println("err connecting to  database", err)
		return nil, err
	}

	//concat('[',ST_X(geom),',',ST_Y(geom), ']') as geom
	sql := "select fid,osm_id, name, ST_AsGeoJson(geom) as geom from view_place_routing"
	stmt, err := conn.Prepare(sql)
	if err != nil {
		fmt.Println("err connecting to  database", err)

		return nil, err
	}
	//source=1529&target=2225
	//(source, target

	rows, err := stmt.Query() // (source, target)
	if err != nil {
		fmt.Println("err connecting to  database", err)
		return nil, err
	}

	// rows, err := conn.Query("select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost, cost as reverse_cost FROM public.roads', $1 , $2 ,true)")
	// checkErr(err)

	var pois []POI
	err = scan.Rows(&pois, rows)
	util.CheckErr(err)
	// fmt.Printf("%#v", steps)

	features := make([]Feature, len(pois))

	for i := 0; i < len(pois); i++ {
		var geojson GeoJSON
		b := []byte(pois[i].GeomJSON)
		json.Unmarshal(b, &geojson)
		//pois[i].Geom = geojson

		f := Feature{}
		f.Type = "Feature"
		f.Geom = geojson
		f.POI = &pois[i]
		features[i] = f
	}

	rows.Close()
	stmt.Close()

	col := FeatureCollection{}
	col.Type = "FeatureCollection"
	col.Feature = &features
	col.Name = "POIs"

	//col.POI = &pois
	return &col, nil
}

func FindNearestNodeForPOI(id int) (*Node, error) {

	conn, err := pgsql.GetConnection()
	if err != nil {
		fmt.Println("err connecting to  database", err)
		return nil, err
	}

	sql := fmt.Sprintf("select v.id, ST_X(v.the_geom) as lon,ST_Y(v.the_geom) as lat from roads_vertices_pgr v,view_place_routing p where v.ein > 0 and v.eout > 0 and p.osm_id=%d order by v.the_geom <-> p.geom limit 1", id)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		fmt.Println("err connecting to  database", err)

		return nil, err
	}
	//source=1529&target=2225
	//(source, target

	rows, err := stmt.Query() // (source, target)
	if err != nil {
		fmt.Println("err connecting to  database", err)
		return nil, err
	}

	// rows, err := conn.Query("select edge FROM pgr_dijkstra('SELECT osm_id as id, source, target, cost, cost as reverse_cost FROM public.roads', $1 , $2 ,true)")
	// checkErr(err)

	var node Node
	err = scan.Row(&node, rows)
	util.CheckErr(err)

	return &node, err

}
