package server

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestType int

const (
	UnitTest = iota + 1
	IntegrationTest
	_
	SmokeTest
	LoadTest
	PerformanceTest
	AcceptenceTest
)

type Test struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description" validate:"required"`
	TestType    TestType           `json:"type" bson:"type" validate:"gte=1,lte10"` //validate:"validateTestType"
	RunAt       time.Time          `json:"run" bson:"run"`
	// ... `json:"-"`
}

func (p *Test) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// Products is a collection of Product
type Tests []*Test

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func (p *Tests) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Test) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Test) Validate() error {
	validate := validator.New()
	// validate.RegisterValidation("testtype", validateTestType)

	return validate.Struct(p)
}

// func validateSKU(fl validator.FieldLevel) bool {
// 	// sku is of format abc-absd-dfsdf
// 	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
// 	matches := re.FindAllString(fl.Field().String(), -1)

// 	if len(matches) != 1 {
// 		return false
// 	}

// 	return true
// }

// GetProducts returns a list of products
func GetTests() Tests {
	return productList
}

func GetTestById(id primitive.ObjectID) (*Test, error) {
	t, _, e := findTestById(id)

	return t, e
}

func AddTest(p *Test) error {
	p.ID = primitive.NewObjectID()
	productList = append(productList, p)
	return nil
}

func UpdateTest(id primitive.ObjectID, p *Test) error {
	_, pos, err := findTestById(id)
	if err != nil {
		return err
	}

	p.ID = primitive.NewObjectID()
	productList[pos] = p

	return nil
}

func DeleteTest(id primitive.ObjectID) error {
	_, pos, err := findTestById(id)
	if err != nil {
		return err
	}

	productList = removeIndexInSlice(productList, pos)

	return nil
}

func removeIndexInSlice(slice []*Test, i int) []*Test {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

var ErrTestNotFound = fmt.Errorf("Test not found")

func findTestById(id primitive.ObjectID) (*Test, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrTestNotFound
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Test{
	&Test{
		ID:          primitive.NewObjectID(),
		Name:        "Latte",
		Description: "Frothy milky coffee",
	},
	&Test{
		ID:          primitive.NewObjectID(),
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
	},
}
