package routing

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PositionEvent struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Position Position           `json:"pos" bson:"pos"`
	Time     time.Time          `json:"time" bson:"time"`
}

type RouteEvent struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Route Route              `json:"route" bson:"route"`
}

type TruckActionEvent struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	ActionType ActionType         `json:"action" bson:"action"`
}

