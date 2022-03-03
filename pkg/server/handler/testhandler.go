package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	data "github.com/andidroid/testgo/pkg/server/data"
	// "data"
)

// Tests is a http.Handler
type TestHandler struct {
	l *log.Logger
}

// NewTestHandler creates a products handler with the given logger
func NewTestHandler(l *log.Logger) *TestHandler {
	return &TestHandler{l}
}

// ServeHTTP is the main entry point for the handler and satisfies the http.Handler
// interface
func (p *TestHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		p.HandleGetRequest(rw, r)
	} else if r.Method == http.MethodPost {
		p.HandlePostRequest(rw, r)
	} else if r.Method == http.MethodPut {
		p.HandlePutRequest(rw, r)
	} else if r.Method == http.MethodDelete {
		p.HandleDeleteRequest(rw, r)
	} else {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (p *TestHandler) getId(r *http.Request) (int, error) {
	p.l.Println("getId - find id for request:", r.Method, r.URL.Path)
	// expect the id in the URI
	reg := regexp.MustCompile(`/([0-9]+)`)
	g := reg.FindAllStringSubmatch(r.URL.Path, -1)
	p.l.Println(g)
	if len(g) != 1 {
		p.l.Println("Invalid URI more than one id")
		return -1, fmt.Errorf("Invalid URI more than one id", r.URL.Path)
	}

	if len(g[0]) != 2 {
		p.l.Println("Invalid URI more than one capture group")
		return -1, fmt.Errorf("Invalid URI more than one id", r.URL.Path)
	}

	idString := g[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		p.l.Println("Invalid URI unable to convert to numer", idString)
		return -1, fmt.Errorf("Invalid URI more than one id", r.URL.Path)
	}
	return id, nil
}

// getTestHandler returns the products from the data store
func (p *TestHandler) HandleGetRequest(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Tests")

	id, err := p.getId(r)
	// if err != nil {
	// 	http.Error(rw, "Invalid URI", http.StatusBadRequest)
	// 	return
	// }

	//id := strings.TrimPrefix(r.URL.Path, "/test/")

	//if id == "" {
	if id == -1 {
		// fetch the products from the datastore
		lp := data.GetTests()
		// serialize the list to JSON
		err := lp.ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
	} else {
		//id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(rw, "Unable to convert id", http.StatusBadRequest)
			return
		}
		lp, err := data.GetTestById(id)
		if err == data.ErrTestNotFound {
			http.Error(rw, "Test not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(rw, "Test not found", http.StatusInternalServerError)
		}
		err = lp.ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
	}

}

func (p *TestHandler) HandlePostRequest(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Test")

	prod := r.Context().Value(KeyTest{}).(data.Test)
	err := data.AddTest(&prod)
	if err != nil {
		http.Error(rw, "Invalid Test", http.StatusBadRequest)
		return
	}
}

func (p *TestHandler) HandlePutRequest(rw http.ResponseWriter, r *http.Request) {
	id, err := p.getId(r)
	if err != nil {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Test", id)
	prod := r.Context().Value(KeyTest{}).(data.Test)

	err = data.UpdateTest(id, &prod)
	if err == data.ErrTestNotFound {
		http.Error(rw, "Test not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Test not found", http.StatusInternalServerError)
		return
	}
}

func (p *TestHandler) HandleDeleteRequest(rw http.ResponseWriter, r *http.Request) {
	id, err := p.getId(r)
	if err != nil {
		http.Error(rw, "Invalid URI", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle DELETE Test", id)

	err = data.DeleteTest(id)
	if err == data.ErrTestNotFound {
		http.Error(rw, "Test not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Test not found", http.StatusInternalServerError)
		return
	}
}

type KeyTest struct{}

func (p *TestHandler) MiddlewareValidateTest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Test{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Error reading product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyTest{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
