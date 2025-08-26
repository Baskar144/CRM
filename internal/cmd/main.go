package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type customers struct {
	Id        uint
	Name      string
	Role      string
	Email     string
	Phone     uint64
	Contacted bool
}

var sampleDetails []customers

func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "./../static/about.html")
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sampleDetails)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//to fetch the path variable
	params := mux.Vars(r)

	idparam, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	for _, v := range sampleDetails {
		if int(v.Id) == idparam {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(v)
		}
	}
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCustomerDetail customers //slice to be defined if multiple customers to be added.
	request_body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(request_body, &newCustomerDetail)
	sampleDetails = append(sampleDetails, newCustomerDetail)
	//sampleDetails = append(sampleDetails, newCustomerDetail...) // to add multiple customers
	json.NewEncoder(w).Encode(sampleDetails)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var latestDetail customers
	request_body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(request_body, &latestDetail)

	params := mux.Vars(r)

	idparams, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	for k, v := range sampleDetails {
		if int(v.Id) == idparams {
			sampleDetails = append(sampleDetails[:k], sampleDetails[k+1:]...)
			sampleDetails = append(sampleDetails, latestDetail)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(sampleDetails)
		}
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	idparams, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	for k, v := range sampleDetails {
		if int(v.Id) == idparams {
			sampleDetails = append(sampleDetails[:k], sampleDetails[k+1:]...)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Customer deleted successfully",
			})
		}
	}
}

func main() {
	fmt.Println("Welcome to the backend service of CRM tool!!!")
	sampleDetails = []customers{
		{1, "Tom", "Developer", "tom1@gmail.com", 9876543210, true},
		{2, "David", "QA Engineer", "david2@gmail.com", 9876543210, false},
		{3, "Stuart", "Project manager", "stuart3@gmail.com", 9876543210, true},
	}

	fmt.Println("Please find the customer details below:")

	for _, detail := range sampleDetails {
		fmt.Println(detail)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", about).Methods("GET")
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("Server started...")

	http.ListenAndServe(":3000", router)
}
