package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	// "io/ioutil"
	"math/rand"
	"net/http"
)

type Customer struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
	Contacted bool   `json:"contacted"`
}

// Customer ID is also used as key in database map.
var database = map[string]Customer{
	"023004163": {
		"023004163",
		"Alan Grant",
		"Paleontologist",
		"agrant@du.edu",
		3038898552,
		true,
	},
	"490520604": {
		"490520604",
		"Ellie Sattler",
		"Paleobotanist",
		"esattler@du.edu",
		3038281386,
		true,
	},
	"344093830": {
		"344093830",
		"Ian Malcolm",
		"Mathematician",
		"imalcolm@math.utexas.edu",
		5128416655,
		true,
	},
	"869930202": {
		"869930202",
		"Donald Gennaro",
		"Attorney",
		"dgennaro@cowanswainross.com",
		4158845018,
		true,
	},
	"400025134": {
		"400025134",
		"Lex Murphy",
		"Student",
		"lex911@aol.com",
		7186177299,
		false,
	},
	"730856990": {
		"730856990",
		"Tim Murphy",
		"Student",
		"tim921@aol.com",
		7186177299,
		false,
	},
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(nil)
	}
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(nil)
	} else {
		if customer.Id == "" {
			// If the request did not include an id, then we
			// randomly generate one. We keep trying until
			// there is no conflict with an existing id or
			// until we have tried and failed 1000 times.
			// This simple approach is for example purposes.
			// In reality, we would probably do something
			// more sophisticated to generate unique ids.
			random := rand.New(rand.NewSource(0))
			failureLimit := int(1e3)
			failureCount := 0
			for failureCount < failureLimit {
				idInt := random.Intn(int(1e9))
				idString := fmt.Sprintf("%09d", idInt)
				_, conflict := database[idString]
				if !conflict {
					customer.Id = idString
					break
				}
				failureCount++
			}
		}
		stillNoId := customer.Id == ""
		_, conflict := database[customer.Id]
		if stillNoId || conflict {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(nil)
		} else {
			database[customer.Id] = customer
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(database)
		}
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	var updatedInfo Customer
	err := json.NewDecoder(r.Body).Decode(&updatedInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(nil)
	} else if _, ok := database[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
	} else {
		if updatedInfo.Id == "" {
			// Since the id can be missing from the data
			// when adding a customer, we also allow it to
			// be missing when updating a customer.
			// When updating, we can infer it from the url.
			updatedInfo.Id = id
		}
		idChanged := updatedInfo.Id != id
		_, updatedIdExists := database[updatedInfo.Id]
		if idChanged && updatedIdExists {
			// Do not allow a customer's id to be updated
			// to the same id as another customer.
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(nil)
		} else {
			database[updatedInfo.Id] = updatedInfo
			if idChanged {
				// If the update changes the customer's id,
				// then a new entry is created, in which case
				// we must remember to delete the old one.
				delete(database, id)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedInfo)
		}
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if _, ok := database[id]; ok {
		delete(database, id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(database)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	port := "3000"
	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, router)
}
