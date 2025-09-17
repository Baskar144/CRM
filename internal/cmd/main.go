package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Customer struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

var db *gorm.DB

func connectDatabase() {
	var err error
	//dsn - data source name
	dsn := "host=localhost user=postgres password=Test@123 dbname=customersdb1 port=5432 sslmode=disable"
	if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Database connection successful")
}

func migrateDatabase() {
	if err := db.AutoMigrate(&Customer{}); err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	log.Println("Database migrated successfully")
}

func seedDetailsToDatabase() {
	var count int64

	if err := db.Model(&Customer{}).Count(&count).Error; err != nil {
		log.Fatalf("Error checking customer count: %v", err)
	}

	if count == 0 {
		customersList := []Customer{
			{Id: uuid.NewString(), Name: "Tom", Role: "Developer", CreatedAt: time.Now().Format(time.RFC3339)},
			{Id: uuid.NewString(), Name: "David", Role: "QA Engineer", CreatedAt: time.Now().Format(time.RFC3339)},
			{Id: uuid.NewString(), Name: "Stuart", Role: "Project manager", CreatedAt: time.Now().Format(time.RFC3339)},
		}

		for _, customer := range customersList {
			if err := db.Create(&customer).Error; err != nil {
				log.Fatalf("Error creating a record in the database: %v", err)
			}
			log.Println("Record created successfully")
		}
	} else {
		log.Printf("Database already contains %d records", count)
	}
}

func main() {

	connectDatabase()
	migrateDatabase()
	seedDetailsToDatabase()

	//Defining the gin router with the default middleware
	router := gin.Default()

	// Registering the welcome route
	router.GET("/welcome", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to CRM tool using Gin framework!")
	})

	// Registering the route to display the customers list
	router.GET("/customers", func(ctx *gin.Context) {
		var latestDetails []Customer

		if err := db.Find(&latestDetails).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching the details from the database"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"customers": latestDetails})
	})

	// Registering the route to fetch a single customer by its ID
	router.GET("/customer/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var customer Customer

		if err := db.Where("id = ?", id).First(&customer).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer record"})
			}
		}
		ctx.JSON(http.StatusOK, customer)
	})

	// Registering the route to create a new customer
	router.POST("/customer", func(ctx *gin.Context) {
		var customer Customer

		if err := ctx.ShouldBindJSON(&customer); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		customer.Id = uuid.NewString()
		customer.CreatedAt = time.Now().Format(time.RFC3339)

		if err := db.Create(&customer).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
			return
		}
		ctx.JSON(http.StatusCreated, customer)
	})

	//Registering routes to update a specific customer record
	router.PATCH("/customer/:id", func(ctx *gin.Context) {
		var customer Customer
		id := ctx.Param("id")

		if err := ctx.ShouldBindJSON(&customer); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		customer.Id = id

		if err := db.Model(&Customer{}).Where("id = ?", id).Updates(customer).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the record"})
			return
		}
		ctx.JSON(http.StatusAccepted, customer)
	})

	// Registering a route to delete a customer record
	router.DELETE("/customer/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		if err := db.Delete(&Customer{}, "id = ?", id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the record"})
			return
		}
		ctx.JSON(http.StatusAccepted, gin.H{"message": "customer record deleted successfully"})
	})

	log.Println("Server running at port 3000")

	// Starting the server
	router.Run(":3000")
}
