package routing

type Position struct {
	Lon float64 `json:"lon" bson:"lon"`
	Lat float64 `json:"lat" bson:"lat"`
}
