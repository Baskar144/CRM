# Customer Relationship Management

This project focuses on building a simple backend service for Customer Relationship Management (CRM). The project uses GORM, net/http, and Postgres (SQL database to store the records in database) standard libraries from Go.

# Descrition

As part of CRM, the below operations will be performed,

1. Get a single/ multiple customer details
2. Add a new customer details
3. Update the existing customer details
4. Delete a customer record

## Getting Started

1. Install the latest version of Go
2. Install an IDE to view the code repository (preferrably VS Code)
3. Install Postman to test the APIs

## Executing Program

1. Install the below third party libraries to define and handle the http requests.
    a. `gorilla/mux`
    b. `gorm.io/gorm`
    c. `gorm.io/driver/postgres`
2. `go run main.go` - command used to run and execute the program.