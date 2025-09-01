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

func getId() uint {
	var highestId uint = 0

	for _, customer := range sampleDetails {
		if customer.Id > uint(highestId) {
			highestId = customer.Id
		}
	}
	return highestId + 1
}

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
	idparam, _ := strconv.Atoi(params["id"])

	for _, v := range sampleDetails {
		if int(v.Id) == idparam {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(v)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "ID not found"}`))
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newCustomerDetail customers //slice to be defined if multiple customers to be added.
	request_body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(request_body, &newCustomerDetail)

	newCustomerDetail.Id = getId()
	sampleDetails = append(sampleDetails, newCustomerDetail)
	//sampleDetails = append(sampleDetails, newCustomerDetail...) // to add multiple customers
	json.NewEncoder(w).Encode(sampleDetails)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var latestDetail customers
	request_body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(request_body, &latestDetail)

	// json.NewDecoder(r.Body).Decode(&latestDetail) - alternate to decode the JSON request to Go data structure

	params := mux.Vars(r)
	idparams, _ := strconv.Atoi(params["id"])

	for k, v := range sampleDetails {
		if int(v.Id) == idparams {
			v = latestDetail
			sampleDetails[k] = v
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Customer updated successfully",
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "ID not found"}`))
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	idparams, _ := strconv.Atoi(params["id"])

	for k, v := range sampleDetails {
		if int(v.Id) == idparams {
			sampleDetails = append(sampleDetails[:k], sampleDetails[k+1:]...)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Customer deleted successfully",
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "ID not found"}`))
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
