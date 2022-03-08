package handler

import (
	"fmt"
	"io"
	"log"
	"time"

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
		for {
			time.Sleep(time.Second * 10)
			now := time.Now().Format("2006-01-02 15:04:05")
			currentTime := fmt.Sprintf("The Current Time Is %v", now)

			// Send current time to clients message channel
			stream.Message <- currentTime
		}
	}()

	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		if msg, ok := <-stream.Message; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

	// //
}

//

type EventStreamHandler struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan string

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
		Message:       make(chan string),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
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
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
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
