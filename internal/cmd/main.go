package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Customer struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     uint64 `json:"phone"`
	Contacted bool   `json:"contacted"`
}

var db *gorm.DB
var customersList []Customer

func initializeDatabase() {
	dsn := "host=localhost user=postgres password=Test@123 dbname=customersdb port=5432 sslmode=disable"

	var err error
	if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = db.AutoMigrate(&Customer{}); err != nil {
		log.Fatalf("Failed to migrate the database schema: %v", err)
	}

	customersList = []Customer{
		{Name: "Tom", Role: "Developer", Email: "tom1@gmail.com", Phone: 9876543210, Contacted: true},
		{Name: "David", Role: "QA Engineer", Email: "david2@gmail.com", Phone: 9876543210, Contacted: false},
		{Name: "Stuart", Role: "Project manager", Email: "stuart3@gmail.com", Phone: 9876543210, Contacted: true},
	}

	for i, customer := range customersList {
		var exisingCustomer Customer
		result := db.First(&exisingCustomer, i+1)
		if result.RowsAffected == 0 {
			db.Create(&customer)
		}
	}

	fmt.Println("Database connection established, schema migrated, and initial data seeded!")
}

func about(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid http method found", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, "./../static/about.html")
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid http method found", http.StatusMethodNotAllowed)
		return
	}

	var customers []Customer
	if err := db.Find(&customers).Error; err != nil {
		http.Error(w, "Error fetching the records from the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func getCustomerById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid http method found", http.StatusMethodNotAllowed)
		return
	}
	//to fetch the path variable
	param, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid customer ID", http.StatusInternalServerError)
		return
	}

	var customer Customer
	result := db.First(&customer, param)
	if result.Error != nil {
		http.Error(w, `{"error": "Customer details not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid http method found", http.StatusMethodNotAllowed)
		return
	}

	var newCustomerDetail Customer
	request_body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(request_body, &newCustomerDetail); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := db.Create(&newCustomerDetail).Error; err != nil {
		http.Error(w, "Cannot insert data into database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomerDetail)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "invalid http method found", http.StatusMethodNotAllowed)
		return
	}

	var latestDetail Customer
	json.NewDecoder(r.Body).Decode(&latestDetail) //alternate to decode the JSON request to Go data structure

	params := mux.Vars(r)
	idparams, _ := strconv.Atoi(params["id"])

	for k, v := range customersList {
		if int(v.Id) == idparams {
			v = latestDetail
			customersList[k] = v
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Customer updated successfully",
			})
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "ID not found"}`))
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "invalid http method found", http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)
	idparams, _ := strconv.Atoi(params["id"])

	for k, v := range customersList {
		if int(v.Id) == idparams {
			customersList = append(customersList[:k], customersList[k+1:]...)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Customer deleted successfully",
			})
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "ID not found"}`))
}

func main() {
	fmt.Println("Welcome to the backend service of CRM tool!!!")
	initializeDatabase()
	fmt.Println("Please find the customer details below:")
	fmt.Println(customersList)

	router := mux.NewRouter()

	router.HandleFunc("/", about).Methods("GET")
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customer/{id}", getCustomerById).Methods("GET")
	router.HandleFunc("/customer", createCustomer).Methods("POST")
	router.HandleFunc("/customer/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customer/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("Server started...")

	http.ListenAndServe(":3000", router)
}
