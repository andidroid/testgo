package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/andidroid/testgo/pkg/redis"
	"github.com/andidroid/testgo/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID         primitive.ObjectID `json:"id" bson:"_id" db:"id"`
	Name       string             `json:"name" bson:"name" db:"name"`
	Position   *Position          `json:"pos" bson:"-" db:"-"`
	Route      *Route             `json:"route" bson:"-" db:"-"`
	Base       *Base              `json:"base" bson:"base" db:"base"`
	ActionType ActionType         `json:"action" bson:"action" db:"-"`
}

type Base struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" db:"id"`
	Name string             `json:"name" bson:"name" db:"name"`
	Geom Geometry           `json:"geometry" bson:"geom"`
}

func NewTruck(name string) *Truck {
	base := Base{Name: "Start", Geom: Geometry{Type: "Point", Coordinates: []float64{11.50, 52.00}}}
	position := Position{Lon: base.Geom.Coordinates[0], Lat: base.Geom.Coordinates[1]}

	return &Truck{ID: primitive.NewObjectID(), Name: name, Base: &base, ActionType: STOPPED, Position: &position}
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
		ID:    primitive.NewObjectID(),
		Name:  truck.Name,
		Route: *truck.Route,
	}
	payload, err := json.Marshal(re)
	util.CheckErr(err)
	redis.GetClient().Publish(context.Background(), "route", payload)
}

func (truck *Truck) sendPositionEvent() {
	pe := PositionEvent{
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
		client := GetClient()
		url := fmt.Sprintf("%s/node/%d", ROUTING_SERVCICE_URL, nodeId)
		resp, err := client.Get(url)
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
	fmt.Println("Order: ", order)
	order.StartTime = time.Now()

	tsp := order.Tsp
	// TODO: order.Status=RUNNING

	fmt.Println("Start: ", tsp.Start)

	fmt.Println("TSP: ", tsp)

	truck.SetActionType(PREPARING)
	fmt.Println(truck)

	var nearestCurrentNode Node

	client := GetClient()
	url := fmt.Sprintf("%s/node/source?lon=%f&lat=%f", ROUTING_SERVCICE_URL, truck.Position.Lon, truck.Position.Lat)
	fmt.Println(url)
	resp, err := client.Get(url)
	util.CheckErr(err)
	fmt.Println("server GET Node: ", url, resp, err)

	nodeBytes, err := ioutil.ReadAll(resp.Body)
	util.CheckErr(err)
	err = json.Unmarshal(nodeBytes, &nearestCurrentNode)

	util.CheckErr(err)
	fmt.Println("respond Node: ", nearestCurrentNode)

	truck.SetActionType(OUTWARD)
	truck.runRoute(nearestCurrentNode.ID, tsp.Steps[0].Node)

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
		//
		//fmt.Printf("coord %d: lon=%f, lat=%f", i, coord[0], coord[1])
		//PositionEvent
		//fmt.Println(truck)
		//time.Sleep(1 * time.Second)
		err := truck.runRoute(source, target)
		util.CheckErr(err)
		if err != nil {
			return
		}

		//{"type":"LineString","coordinates":[[11.6072681,51.7953606],
		truck.SetActionType(RESTING)
		time.Sleep(5 * time.Second)
		//TODO loop
		//return
	}

	fmt.Printf("--- start order %s %s ---", truck.Name, order.Name)

}

func (truck *Truck) runRoute(source int64, target int64) error {
	route, err := truck.requestRoute(source, target)
	util.CheckErr(err)
	if err != nil {
		return err
	}

	truck.SetRoute(&route)

	truck.SetActionType(RUNNING)

	for i := 0; i < len(route.Coordinates); i++ {
		coord := route.Coordinates[i]

		truck.Position = &Position{Lon: coord[0], Lat: coord[1]}

		truck.sendPositionEvent()

		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (truck *Truck) requestRoute(source int64, target int64) (Route, error) {

	fmt.Println("truck.requestRoute: ", source, target)
	routeKey := fmt.Sprintf("route_geometry_%d_%d", source, target)

	var route Route
	routeString, err := redis.GetClient().Get(context.Background(), routeKey).Result()
	util.CheckErr(err)

	fmt.Println("redis GET Route: ", routeKey, routeString, err)
	var routeBytes []byte

	if err != nil || routeString == "" {
		//request from server

		client := GetClient()
		url := fmt.Sprintf("%s/routing/geometry?source=%d&target=%d", ROUTING_SERVCICE_URL, source, target)
		resp, err := client.Get(url)
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
