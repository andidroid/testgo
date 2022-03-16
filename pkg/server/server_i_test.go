package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	data "github.com/andidroid/testgo/pkg/server/data"
	"github.com/stretchr/testify/assert"
)

// func SetupRouter() *gin.Engine {
// 	router := gin.Default()
// 	return router
// }

var itest string

func init() {
	flag.StringVar(&itest, "itest", "", "the foo bar bang")
}

func TestListRecipesHandler(t *testing.T) {
	if itest != "it" {
		fmt.Print("skip itest")
		t.Skip("Skipping for itest != it")
	}
	// r := SetupRouter()
	//r.GET("/tests", testHandler.HandleGetAllTestsRequest)

	r := CreateRouter()
	req, _ := http.NewRequest("GET", "/tests", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var tests []data.Test
	json.Unmarshal([]byte(w.Body.String()), &tests)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 9, len(tests))
}

// // TestHelloName calls greetings.Hello with a name, checking
// // for a valid return value.
// func TestHelloRequest(t *testing.T) {

// 	resp, err := http.Get("http://localhost:8090/hello")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//Convert the body to type string
// 	sb := string(body)
// 	log.Printf(sb)

// 	want := regexp.MustCompile("hello\n")
// 	if !want.MatchString(sb) || err != nil {
// 		t.Fatalf(`Hello("Gladys") = %q, want match for %#q, nil`, sb, want)
// 	}
// }
