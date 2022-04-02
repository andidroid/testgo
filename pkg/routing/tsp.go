package routing

import (
	"fmt"

	"github.com/andidroid/testgo/pkg/pgsql"
	"github.com/blockloop/scan"
)

type Tsp struct {
	Start  int       `json:"start" bson:"start"`
	Steps  []TspStep `json:"steps" bson:"steps"`
	Places string    `json:"places" bson:"places"`
}

type TspStep struct {
	Seq     int     `json:"seq" bson:"seq"`
	Node    int64   `json:"node" bson:"node"`
	Cost    float64 `json:"cost" bson:"cost"`
	AggCost float64 `json:"agg_cost" bson:"agg_cost"`
}

//ids []int
func CalculateTSP(start int, ids string) (*Tsp, error) {

	conn, err := pgsql.GetConnection()
	if err != nil {
		fmt.Println("err connecting to  database", err)
		return nil, err
	}

	cost := func() string {
		if DIRECTED_GRAPH {
			return "reverse_cost"
		} else {
			return "cost"
		}
	}()

	var inQuery string
	//ids != nil &&
	if len(ids) > 0 {
		//strings.Join(ids, ",")
		inQuery = fmt.Sprintf("WHERE view_place_routing.osm_id IN(%s)", ids)
	} else {
		inQuery = ""
	}

	sql := fmt.Sprintf("SELECT * FROM pgr_TSP($$ SELECT * FROM pgr_dijkstraCostMatrix('SELECT osm_id as id, source, target, cost, %s FROM view_routing',( with nearestVertices as ( SELECT a.id from view_place_routing , lateral ( select id, the_geom from roads_vertices_pgr where ein > 0 and eout > 0 order by roads_vertices_pgr.the_geom <-> view_place_routing.geom limit 1 ) a %s) select array_agg(id) from nearestVertices), directed := %t)$$, start_id := %d)", cost, inQuery, DIRECTED_GRAPH, start)
	fmt.Println(sql)
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

	var steps []TspStep
	err = scan.Rows(&steps, rows)

	// fmt.Printf("%#v", steps)

	var tsp = Tsp{Start: start, Steps: steps, Places: ids}
	rows.Close()
	stmt.Close()

	return &tsp, nil
}
