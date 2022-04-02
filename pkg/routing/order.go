package routing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andidroid/testgo/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type OrderState int

const (
	ORDER_PLANNED = iota + 1
	ORDER_PROCESSING
	ORDER_STARTED
	ORDER_FINISHED
)

type Order struct {
	ID         primitive.ObjectID `json:"id" bson:"_id" db:"id"`
	Name       string             `json:"name" bson:"name" db:"name"`
	OrderState OrderState         `json:"state" bson:"state" db:"-"`
	StartTime  time.Time          `json:"starttime" bson:"starttime" db:"starttime"`
	EndTime    time.Time          `json:"endtime" bson:"endtime" db:"endtime"`
	Tsp        *Tsp               `json:"-" bson:"-" db:"-"`
	TruckID    primitive.ObjectID `json:"truckid" bson:"_truckid" db:"truckid"`
}

func NewOrder(name string) *Order {
	//, tsp *Tsp, truckID primitive.ObjectID
	//, Tsp: tsp, TruckID: truckID
	return &Order{ID: primitive.NewObjectID(), Name: name, OrderState: ORDER_PLANNED}
}

func (order *Order) RequestTSP(startNodeId int) error {

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	//TODO add ids to request
	url := fmt.Sprintf("http://localhost/routing/tsp?start=%d", startNodeId)
	resp, err := client.Get(url)
	util.CheckErr(err)
	// fmt.Printf("calling /routing/tsp?start=84845: %s", resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
		return fmt.Errorf("response error %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	util.CheckErr(err)

	var tsp *Tsp
	err = json.Unmarshal(body, &tsp)
	util.CheckErr(err)
	fmt.Println("read TSP: ", tsp)
	order.Tsp = tsp
	return nil
}

func (order *Order) AssignTruck(truckID primitive.ObjectID) {
	order.TruckID = truckID
}
