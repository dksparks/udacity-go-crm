package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	// "io"
	"math/rand"
	"net/http"
)

type Customer struct {
	Name  string `json:"name"`
	Adult bool   `json:"adult"`
}

// The customer database uses nine-digit ids as keys for simplicity.
// In reality, we would probably do something more sophisticated.
var database = map[string]Customer{
	"023004163": Customer{"Alan Grant", true},
	"490520604": Customer{"Ellie Sattler", true},
	"344093830": Customer{"Ian Malcolm", true},
	"869930202": Customer{"Donald Gennaro", true},
	"400025134": Customer{"Lex Murphy", false},
	"730856990": Customer{"Tim Murphy", false},
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(database)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if customer, ok := database[id]; ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(customer)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
    fmt.Println(customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// Generate a new database id for the customer,
		// and keep trying until there is no conflict or
		// until we have tried and failed 1000 times.
		// This simple approach is for example purposes.
		// In reality, we would probably do something
        // more sophisticated.
		random := rand.New(rand.NewSource(0))
		failureLimit := int(1e3)
		failureCount := 0
		for failureCount < failureLimit {
			idNum := random.Intn(int(1e9))
			id := fmt.Sprintf("%09d", idNum)
			_, conflict := database[id]
			if !conflict {
				database[id] = customer
				w.WriteHeader(http.StatusCreated)
				break
			}
			failureCount++
		}
		if failureCount == failureLimit {
			w.WriteHeader(http.StatusConflict)
		}
		json.NewEncoder(w).Encode(database)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	port := "3000"
	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, router)
}
