package routing

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Place struct {
	ID          primitive.ObjectID `json:"id" bson:"_id" db:"id"`
	Name        string             `json:"name" bson:"name" db:"name"`
	Geometry    Geometry           `json:"geom" bson:"geom" db:"geom"`
	OsmId       int64              `json:"osm_id" bson:"osm_id" db:"osm_id"`
	Description string             `json:"description" bson:"description" db:"description"`
}

func NewPlace(name string) *Place {
	//, tsp *Tsp, truckID primitive.ObjectID
	//, Tsp: tsp, TruckID: truckID
	return &Place{ID: primitive.NewObjectID(), Name: name}
}
