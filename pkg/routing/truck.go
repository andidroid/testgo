package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andidroid/testgo/pkg/redis"
	"github.com/andidroid/testgo/pkg/util"
)

type ActionType int

const (
	STOPPED = iota + 1
	RUNNING
	PAUSED
	PREPARING
	_
	STARTED
	FINISHED
	HOMEWARD
	OUTWARD
	RESTING
)

type Truck struct {
	ID         int64      `json:"id" bson:"_id" db:"id"`
	Name       string     `json:"name" bson:"name" db:"name"`
	Position   *Position  `json:"pos" bson:"-" db:"-"`
	Route      *Route     `json:"route" bson:"-" db:"-"`
	Base       int64      `json:"base" bson:"base" db:"base"`
	ActionType ActionType `json:"action" bson:"action" db:"-"`
}

func NewTruck(name string) *Truck {
	return &Truck{ID: 1, Name: name, Base: 84845, ActionType: STOPPED}
}

func (truck *Truck) SetActionType(ActionType ActionType) {
	truck.ActionType = ActionType

	truck.sendTruckActionEvent()
}

func (truck *Truck) SetRoute(route *Route) {
	truck.Route = route
	truck.sendRouteEvent()
}

func (truck *Truck) sendTruckActionEvent() {
	pe := TruckActionEvent{
		ID:         truck.ID,
		Name:       truck.Name,
		ActionType: truck.ActionType,
	}
	payload, err := json.Marshal(pe)
	util.CheckErr(err)
	redis.GetClient().Publish(context.Background(), "action", payload)
}

func (truck *Truck) sendRouteEvent() {
	re := RouteEvent{
		ID:    1,
		Name:  truck.Name,
		Route: *truck.Route,
	}
	payload, err := json.Marshal(re)
	util.CheckErr(err)
	redis.GetClient().Publish(context.Background(), "route", payload)
}

func (truck *Truck) sendPositionEvent() {
	pe := PostionEvent{
		ID:       truck.ID,
		Name:     truck.Name,
		Position: *truck.Position,
		Time:     time.Now(),
	}
	payload, err := json.Marshal(pe)
	util.CheckErr(err)
	redis.GetClient().Publish(context.Background(), "position", payload)
}

func (truck *Truck) requestNode(nodeId int64) Node {
	fmt.Println("request Node: ", nodeId)
	nodeKey := fmt.Sprintf("node_%d", nodeId)

	var node Node
	nodeString, err := redis.GetClient().Get(context.Background(), nodeKey).Result()
	util.CheckErr(err)
	fmt.Println("redis GET Node: ", nodeId, nodeString, err)
	var nodeBytes []byte
	if err != nil || nodeString == "" {
		//request from server
		url := fmt.Sprintf("http://localhost/node/%d", nodeId)
		resp, err := http.Get(url)
		util.CheckErr(err)
		fmt.Println("server GET Node: ", nodeId, resp, err)

		nodeBytes, err = ioutil.ReadAll(resp.Body)
		util.CheckErr(err)
		nodeString = string(nodeBytes)
		err = redis.GetClient().Set(context.Background(), nodeKey, nodeString, 5*time.Minute).Err()
		util.CheckErr(err)
		fmt.Println("redis SET Node: ", nodeId, nodeString, err)
	} else {
		nodeBytes = []byte(nodeString)
	}

	//

	err = json.Unmarshal(nodeBytes, &node)
	util.CheckErr(err)
	fmt.Println("respond Node: ", node)
	return node
}

func (truck *Truck) StartOrder(order *Order) {
	fmt.Printf("+++ start order %s %s +++", truck.Name, order.Name)
	tsp := order.Tsp
	// TODO: order.Status=RUNNING
	fmt.Println(tsp.Start)
	truck.SetActionType(STARTED)

	for i := 0; i < len(tsp.Steps)-1; i++ {
		truck.SetActionType(PREPARING)

		source := tsp.Steps[i].Node
		target := tsp.Steps[i+1].Node

		fmt.Println("truck starting route", truck.Name, source, target, i)

		//fmt.Println(resp.Body)
		// fmt.Println(route)
		//
		//
		route, err := truck.requestRoute(source, target, i)
		util.CheckErr(err)
		if err != nil {
			return
		}

		truck.SetRoute(&route)

		//

		truck.SetActionType(RUNNING)

		for i := 0; i < len(route.Coordinates); i++ {
			coord := route.Coordinates[i]
			//fmt.Printf("coord %d: lon=%f, lat=%f", i, coord[0], coord[1])
			truck.Position = &Position{Lon: coord[0], Lat: coord[1]}

			//PositionEvent
			truck.sendPositionEvent()
			//fmt.Println(truck)
			//time.Sleep(1 * time.Second)
			time.Sleep(100 * time.Millisecond)
		}

		//{"type":"LineString","coordinates":[[11.6072681,51.7953606],
		truck.SetActionType(RESTING)
		time.Sleep(5 * time.Second)
		//TODO loop
		//return
	}

	fmt.Printf("--- start order %s %s ---", truck.Name, order.Name)

}

func (truck *Truck) requestRoute(source int64, target int64, id int) (Route, error) {

	fmt.Println("request Route: ", source, target)
	routeKey := fmt.Sprintf("route_%d_%d", source, target)

	var route Route
	routeString, err := redis.GetClient().Get(context.Background(), routeKey).Result()
	util.CheckErr(err)

	fmt.Println("redis GET Route: ", routeKey, routeString, err)
	var routeBytes []byte

	if err != nil || routeString == "" {
		//request from server
		url := fmt.Sprintf("http://localhost/routing/geometry?source=%d&target=%d", source, target)
		resp, err := http.Get(url)
		util.CheckErr(err)
		if err != nil {
			return Route{}, fmt.Errorf("error requesting route: %w", err)
		}
		fmt.Println("server GET Route: ", routeKey, resp, err)

		routeBytes, err = ioutil.ReadAll(resp.Body)
		util.CheckErr(err)
		routeString = string(routeBytes)
		err = redis.GetClient().Set(context.Background(), routeKey, routeString, 5*time.Minute).Err()
		util.CheckErr(err)
		fmt.Println("redis SET Node: ", routeKey, routeString, err)
	} else {
		routeBytes = []byte(routeString)
	}

	err = json.Unmarshal(routeBytes, &route)
	util.CheckErr(err)

	//request start node to check route direction
	node := truck.requestNode(source)

	if node.Lon == route.Coordinates[0][0] && node.Lat == route.Coordinates[0][1] {

		fmt.Println("route is in correct order")
	} else {
		node = truck.requestNode(target)

		if node.Lon == route.Coordinates[0][0] && node.Lat == route.Coordinates[0][1] {

			fmt.Println("route is in wrong order, need reverse order")

			len := len(route.Coordinates)
			for i := 0; i < len/2; i++ {
				j := len - 1 - i
				if i == j {
					break
				}
				ci := route.Coordinates[i]
				cj := route.Coordinates[j]

				route.Coordinates[i] = cj
				route.Coordinates[j] = ci
			}

		} else {
			fmt.Println("error finding correct route order")
		}

	}
	return route, nil
}

type Route struct {
	ID          int64       `json:"id" bson:"_id"`
	Name        string      `json:"name" bson:"name"`
	Source      int         `json:"source" bson:"source"`
	Target      int         `json:"target" bson:"target"`
	Type        string      `json:"type" bson:"type"`
	Coordinates [][]float64 `json:"coordinates" bson:"coordinates"`
}
