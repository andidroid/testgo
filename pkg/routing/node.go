package routing

type Node struct {
	ID  int64   `json:"id" bson:"_id" db:"id"`
	Lon float64 `json:"lon" bson:"lon" db:"lon"`
	Lat float64 `json:"lat" bson:"lat" db:"lat"`
}
