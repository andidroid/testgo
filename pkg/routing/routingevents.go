package routing

import (
	"time"
)

type PostionEvent struct {
	ID       int64     `json:"id" bson:"_id"`
	Name     string    `json:"name" bson:"name"`
	Position Position  `json:"pos" bson:"pos"`
	Time     time.Time `json:"time" bson:"time"`
}

type RouteEvent struct {
	ID    int64  `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Route Route  `json:"route" bson:"route"`
}

type TruckActionEvent struct {
	ID         int64      `json:"id" bson:"_id"`
	Name       string     `json:"name" bson:"name"`
	ActionType ActionType `json:"action" bson:"action"`
}
