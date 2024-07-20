package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Customer struct {
	Name  string `json:"name"`
	Adult bool   `json:"adult"`
}

var database = map[string]Customer{
	"0230": Customer{"Alan Grant", true},
	"4905": Customer{"Ellie Sattler", true},
	"3440": Customer{"Ian Malcolm", true},
	"8699": Customer{"Donald Gennaro", true},
	"4000": Customer{"Lex Murphy", false},
	"7308": Customer{"Tim Murphy", false},
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	port := "3000"
	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, router)
}
