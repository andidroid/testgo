package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloRequest(t *testing.T) {

	resp, err := http.Get("http://localhost:8090/hello")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)

	want := regexp.MustCompile("hello\n")
	if !want.MatchString(sb) || err != nil {
		t.Fatalf(`Hello("Gladys") = %q, want match for %#q, nil`, sb, want)
	}
}
