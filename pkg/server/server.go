package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"

	// "github.com/andidroid/testgo/pkg/server/data"
	"github.com/andidroid/testgo/pkg/mongo"
	"github.com/andidroid/testgo/pkg/server/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var validPath = regexp.MustCompile("^/(hello|health|test|testmongo)/([a-zA-Z0-9]+)$")

func Start() {
	fmt.Println("start server")

	logger := log.New(os.Stdout, "products-api ", log.LstdFlags)

	testHandler := handler.NewTestHandler(logger)

	http.HandleFunc("/hello/", makeHandler(helloHandler))
	http.HandleFunc("/health/", makeHandler(healthHandler))
	http.HandleFunc("/testmongo/", makeHandler(testMongoHandler))
	http.HandleFunc("/test/", testHandler.ServeHTTP)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println(GetLocalIPAddr())
	err := http.ListenAndServe(":8090", nil)
	log.Fatal(err)
}

func GetLocalIPAddr() net.IP {
	// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
	conn, err := net.Dial("udp", "0.0.0.0:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request, title string) {
	w.Header().Set("Content-Type", "text/plain")
	// query := strings.Split(r.URL.RawQuery, "&")
	// name := query[0]
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "hello %s ? %s", title, name)
}

func healthHandler(w http.ResponseWriter, r *http.Request, title string) {
	//fmt.Fprintf(w, "UP")
	//os.Stdout
	enc := json.NewEncoder(w)
	d := map[string]string{"status": "UP", "call": title}
	enc.Encode(d)

	w.Header().Set("Content-Type", "application/json")
	//implicit w.WriteHeader(http.StatusOK)

}

func testMongoHandler(w http.ResponseWriter, r *http.Request, title string) {
	mongo.Main()
	fmt.Fprintf(w, "Test")

	// var t Test

	// // Try to decode the request body into the struct. If there is an error,
	// // respond to the client with the error message and a 400 status code.
	// err := json.NewDecoder(r.Body).Decode(&t)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}
