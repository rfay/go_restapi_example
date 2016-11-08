package main

// Minor improvements to golang json example in
// https://www.thepolyglotdeveloper.com/2016/07/create-a-simple-restful-api-with-golang/
// Mostly, it just uses a map to store people instead of searching through
// an array. There are also some commented-out debugging statements.
// Get the source set up with:
//   go get github.com/rfay/go_restapi_example
// Run it with
//   go run
// Example POST json to create a new person:
//   POST to http://localhost:12345/people/4 with:
//   {"firstname": "ExampleFirstname", "lastname":"NumberFour", "address":{"city":"somewhereIn", "state":"CO"}}

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"
	"github.com/gorilla/mux"
	"os"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people = make(map[string]Person)

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	person, _ := people[params["id"]]

	// person might be empty if it wasn't found in the map
	json.NewEncoder(w).Encode(person)
}

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	//fmt.Fprintf(os.Stderr, "Params: %v\n", params);
	var person Person
	result := json.NewDecoder(req.Body).Decode(&person)
	if result != nil {
		fmt.Fprintf(os.Stderr, "result=%v\n", result)
		return
	}

	//fmt.Fprintf(os.Stderr, "Person: %v\n", person);
	person.ID, _ = params["id"]
	people[person.ID] = person
	json.NewEncoder(w).Encode(people)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, _ := params["id"]

	delete(people, id)
	json.NewEncoder(w).Encode(people)
}

func main() {
	router := mux.NewRouter()
	people["1"] = Person{ID: "1", Firstname: "Nic", Lastname: "Raboy", Address: &Address{City: "Dublin", State: "CA"}}
	people["2"] = Person{ID: "2", Firstname: "Maria", Lastname: "Raboy"}
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
