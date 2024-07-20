package main

import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
)

type Customer struct {
  Name string `json:"name"`
  // note that id is a string, not an int
  Id string `json:"id"`
  Adult bool `json:"adult"`
}

var database = []Customer{
  Customer{"Alan Grant", "0230", true},
  Customer{"Ellie Sattler", "4905", true},
  Customer{"Ian Malcolm", "3440", true},
  Customer{"Donald Gennaro", "8699", true},
  Customer{"Lex Murphy", "4000", false},
  Customer{"Tim Murphy", "7308", false},
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(database)
}

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/customers", getAllCustomers).Methods("GET")
  port := "3000"
  fmt.Println("Server running on port", port)
  http.ListenAndServe(":" + port, router)
}
