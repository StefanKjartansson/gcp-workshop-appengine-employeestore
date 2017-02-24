// +build !appengine

package main

import (
	"log"
	"net/http"

	"github.com/StefanKjartansson/gcp-workshop-appengine-employeestore"
)

func main() {
	store := employees.NewMemoryStore()
	handler := employees.GetHandler(store)
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
