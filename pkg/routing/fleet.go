package routing

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Fleet struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" db:"id"`
	Trucks map[string]*Truck
	Orders map[string]*Order
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

func (fleet Fleet) StartOrder(orderName string, truckName string, startNodeId int) error {

	truck1 := NewTruck(truckName)
	fleet.Trucks[truckName] = truck1

	order := NewOrder(orderName)
	order.RequestTSP(startNodeId)
	fleet.Orders[orderName] = order
	order.AssignTruck(truck1.ID)
	go func() {
		truck1.StartOrder(order)
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
