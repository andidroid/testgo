package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andidroid/testgo/pkg/util"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Fleet struct {
	Trucks map[string]*Truck
	Orders map[string]*Order
}

type Order struct {
	Name  string `json:"name" bson:"name"`
	Truck string `json:"truck" bson:"truck"`
	Tsp   Tsp
}

var fleet Fleet

func init() {
	initFleet()

	//go startRouting()
	//go fleet.StartOrder("order-0815", "trucky", 84845)
}

func GetFleetInstance() *Fleet {
	return &fleet
}

func initFleet() {
	trucks := make(map[string]*Truck)
	orders := make(map[string]*Order)
	fleet = Fleet{Trucks: trucks, Orders: orders}

}

func startRouting() {
	time.Sleep(time.Second * 5)
	truck1 := Truck{ID: 1, Name: "truck1"}
	fmt.Println("Pointer truck 1: ", &truck1)
	fleet.Trucks["truck1"] = &truck1

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("http://localhost/routing/tsp?start=%d", 84845)
	resp, err := client.Get(url)
	util.CheckErr(err)
	// fmt.Printf("calling /routing/tsp?start=84845: %s", resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	util.CheckErr(err)

	var tsp Tsp
	err = json.Unmarshal(body, &tsp)
	util.CheckErr(err)

	order := Order{Name: "order1", Tsp: tsp, Truck: "truck1"}
	fleet.Orders["order1"] = &order

	go func() {

		truck111 := fleet.Trucks["truck1"]
		truck111.StartOrder(&order)
		//fleet.Trucks["truck1"].StartOrder(&order)
	}()
	go func() {
		for {
			time.Sleep(time.Second * 10)

			truck11 := fleet.Trucks["truck1"]
			currPostion := truck11.Position
			fmt.Println("test postion event:", &currPostion, currPostion)
			fmt.Println("Pointer truck 1-1: ", &truck11, truck11)
		}
	}()
}

func (fleet Fleet) StartOrder(orderName string, truckName string, startNodeId int) error {

	truck1 := Truck{ID: 1, Name: truckName}
	fleet.Trucks[truckName] = &truck1

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := fmt.Sprintf("http://localhost/routing/tsp?start=%d", startNodeId)
	resp, err := client.Get(url)
	util.CheckErr(err)
	// fmt.Printf("calling /routing/tsp?start=84845: %s", resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	util.CheckErr(err)

	var tsp Tsp
	err = json.Unmarshal(body, &tsp)
	util.CheckErr(err)

	order := Order{Name: orderName, Tsp: tsp, Truck: truckName}
	fleet.Orders[orderName] = &order

	go func() {
		truck1.StartOrder(&order)
	}()

	return nil
}

func (fleet Fleet) FindFreeTruck() []*Truck {

	var freeTrucks = make([]*Truck, 0, 1)
	for key, truck := range fleet.Trucks {
		fmt.Println("Key:", key, "=>", "Element:", truck)

		truck := fleet.Trucks[key]
		if truck.ActionType == STOPPED || truck.ActionType == HOMEWARD {
			freeTrucks = append(freeTrucks, truck)
		}
	}
	return freeTrucks
}
