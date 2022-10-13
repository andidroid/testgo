package routing

//geom,omitempty

//{"type":"Point","coordinates":[12.4045328,51.7979734]}
type Geometry struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type FeatureCollection struct {
	ID       string     `json:"_id,omitempty" bson:"_id,omitempty"`
	Type     string     `json:"type" bson:"type"`
	Features *[]Feature `json:"features" bson:"features"`
	Name     string     `json:"name" bson:"name"`
}

type Feature struct {
	ID         string      `json:"_id,omitempty" bson:"_id,omitempty"`
	Type       string      `json:"type" bson:"type"`
	Geom       Geometry    `json:"geometry" bson:"geom"`
	Properties interface{} `json:"properties" bson:"properties"`
}

// https://stackoverflow.com/questions/67807153/golang-mongo-geojson
// // Feature Collection
// type FeatureCollection struct {
//     ID       string    `json:"_id,omitempty" bson:"_id,omitempty"`
//     Features []Feature `json:"features" bson:"features"`
//     Type     string    `json:"type" bson:"type"`
// }

// // Individual Feature
// type Feature struct {
//     Type       string     `json:"type" bson:"type"`
//     Properties Properties `json:"properties" bson:"properties"`
//     Geometry   Geometry   `json:"geometry" bson:"geometry"`
// }

// // Feature Properties
// type Properties struct {
//     Name        string `json:"name" bson:"name"`
//     Height      uint64 `json:"height" bson:"height"`
//     Purchased   bool   `json:"purchased" bson:"purchased"`
//     LastUpdated string `json:"last_updated" bson:"last_updated"`
// }

// // Feature Geometry
// type Geometry struct {
//     Type        string      `json:"type" bson:"type"`
//     Coordinates interface{} `json:"coordinates" bson:"coordinates"`
// }
