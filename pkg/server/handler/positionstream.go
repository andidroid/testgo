package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/andidroid/testgo/pkg/redis"
	"github.com/andidroid/testgo/pkg/routing"
	"github.com/andidroid/testgo/pkg/util"
	"github.com/gin-gonic/gin"
)

// type RoutingHandler struct {
// }

// // NewTestHandler creates a products handler with the given logger
// func NewRoutingHandler() *RoutingHandler {
// 	return &RoutingHandler{}
// }

//localhost:8080/routing/stream
func (stream *EventStreamHandler) GetPositionStream(c *gin.Context) {

	go func() {
		ctx := context.Background()
		sub := redis.GetClient().Subscribe(ctx, "position", "route", "action")
		pe := routing.PostionEvent{}
		for {
			//routing.PostionEvent
			msg, err := sub.ReceiveMessage(ctx)
			util.CheckErr(err)
			// fmt.Println(msg)

			switch msg.Channel {
			case "position":
				if err := json.Unmarshal([]byte(msg.Payload), &pe); err != nil {
					util.CheckErr(err)
				} else {
					stream.PositionEventMessage <- pe
				}
			case "route":
				re := routing.RouteEvent{}
				if err := json.Unmarshal([]byte(msg.Payload), &re); err != nil {
					util.CheckErr(err)
				} else {
					stream.RouteEventMessage <- re
				}
			case "action":
				tae := routing.TruckActionEvent{}
				if err := json.Unmarshal([]byte(msg.Payload), &tae); err != nil {
					util.CheckErr(err)
				} else {
					stream.TruckActionEventMessage <- tae
				}
			}

		}
	}()

	/*
		// old implementation
		go func() {

			if true {
				return
			}

			// var route routing.Route
			// route.ID = math.MinInt64

			oldRouteIDs := make(map[string]int64)

			for {
				time.Sleep(time.Second * 1)
				now := time.Now()

				fleet := *routing.GetINSTANCE()

				for key, truck := range fleet.Trucks {
					fmt.Println("Key:", key, "=>", "Element:", truck)

					truck := fleet.Trucks[key]
					if truck.Route == nil || truck.Position == nil {
						// truck not yet started
						continue
					}
					currPostion := truck.Position

					//TODO: loop over all trucks in map

					// pos := routing.Position{
					// 	Lon: float64(rand.Intn(180)),
					// 	Lat: float64(rand.Intn(90)),
					// }

					pos := routing.Position{
						Lon: currPostion.Lon,
						Lat: currPostion.Lat,
					}

					pe := routing.PostionEvent{
						ID:       1,
						Name:     key,
						Position: pos,
						Time:     now,
					}
					fmt.Println("send postion event:", pe)
					// Send current time to clients message channel
					stream.PositionEventMessage <- pe

					//fmt.Println("check route changed: ", route.ID, fleet.Trucks["truck1"].Route.ID)

					oldRouteID, ok := oldRouteIDs[key]
					if !ok {
						oldRouteID = math.MinInt64
					}

					route := truck.Route
					if oldRouteID != route.ID {
						re := routing.RouteEvent{
							ID:    1,
							Name:  key,
							Route: *route,
						}

						oldRouteIDs[key] = route.ID
						stream.RouteEventMessage <- re
					}

				}
			}
		}()
	*/

	// c.Stream(func(w io.Writer) bool {
	// 	// Stream message to client from message channel
	// 	if msg, ok := <-stream.PositionEventMessage; ok {
	// 		c.SSEvent("position", msg)
	// 		return true
	// 	}
	// 	return false
	// })
	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel

		select {
		// Add new available client
		case msg := <-stream.RouteEventMessage:
			c.SSEvent("route", msg)
			return true
		case msg := <-stream.PositionEventMessage:
			c.SSEvent("position", msg)
			return true
		case msg := <-stream.TruckActionEventMessage:
			c.SSEvent("action", msg)
			return true
		}

		return false
	})
	// //
}

//

type EventStreamHandler struct {
	// Events are pushed to this channel by the main events-gathering routine
	PositionEventMessage    chan routing.PostionEvent
	RouteEventMessage       chan routing.RouteEvent
	TruckActionEventMessage chan routing.TruckActionEvent

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

type ClientChan chan string

// Initialize event and Start procnteessing requests
func NewEventStreamHandler() (event *EventStreamHandler) {

	event = &EventStreamHandler{
		PositionEventMessage:    make(chan routing.PostionEvent),
		RouteEventMessage:       make(chan routing.RouteEvent),
		TruckActionEventMessage: make(chan routing.TruckActionEvent),
		NewClients:              make(chan chan string),
		ClosedClients:           make(chan chan string),
		TotalClients:            make(map[chan string]bool),
	}

	go event.listen()

	return event
}

//It Listens all incoming requests from clients.
//Handles addition and removal of clients and broadcast messages to clients.
func (stream *EventStreamHandler) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.PositionEventMessage:
			for clientMessageChan := range stream.TotalClients {
				json, err := json.Marshal(eventMsg)
				if err != nil {
					log.Println(err)
				}
				clientMessageChan <- string(json)
			}
		case eventMsg := <-stream.RouteEventMessage:
			for clientMessageChan := range stream.TotalClients {
				json, err := json.Marshal(eventMsg)
				if err != nil {
					log.Println(err)
				}
				clientMessageChan <- string(json)
			}
		case eventMsg := <-stream.TruckActionEventMessage:
			for clientMessageChan := range stream.TotalClients {
				json, err := json.Marshal(eventMsg)
				if err != nil {
					log.Println(err)
				}
				clientMessageChan <- string(json)
			}
		}

	}
}

func (stream *EventStreamHandler) ServeHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		go func() {
			// Send connection that is closed by client to event server
			<-c.Done()
			stream.ClosedClients <- clientChan
		}()

		c.Next()
	}
}

func HeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
