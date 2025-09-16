package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Customer struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

var customersList []Customer

func seedDetails() {
	customersList = []Customer{
		{Id: uuid.NewString(), Name: "Tom", Role: "Developer", CreatedAt: time.Now()},
		{Id: uuid.NewString(), Name: "David", Role: "QA Engineer", CreatedAt: time.Now()},
		{Id: uuid.NewString(), Name: "Stuart", Role: "Project manager", CreatedAt: time.Now()},
	}
}

func main() {
	seedDetails()
	//Defining the gin router with the default middleware
	router := gin.Default()

	// Registering the welcome route
	router.GET("/welcome", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to CRM tool using Gin framework!")
	})

	// Registering the route to display the customers list
	router.GET("/customers", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"customers": customersList})
	})

	// Registering the route to fetch a single customer by its ID
	router.GET("/customer/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		for _, customer := range customersList {
			if customer.Id == id {
				ctx.JSON(http.StatusOK, customer)
				return
			}
		}
		ctx.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
	})

	// Registering the route to create a new customer
	router.POST("/customer", func(ctx *gin.Context) {
		var customer Customer

		if err := ctx.ShouldBindJSON(&customer); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		customer.Id = uuid.NewString()
		customer.CreatedAt = time.Now()
		customersList = append(customersList, customer)
		ctx.JSON(http.StatusCreated, customer)
	})

	log.Println("Server running at port 3000")

	// Starting the server
	router.Run(":3000")
}
