package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelrknutson/loanprogo/database"
	"github.com/michaelrknutson/loanprogo/models"
	"github.com/michaelrknutson/loanprogo/operation"
	records "github.com/michaelrknutson/loanprogo/record"
	users "github.com/michaelrknutson/loanprogo/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/users", users.GetUsers)
	app.Get("/api/v1/users/:id", users.GetUser)
	app.Post("/api/v1/users/", users.CreateUser)
	app.Post("/api/v1/users/login", users.LoginUser)
	app.Put("/api/v1/users/:id", users.UpdateUser)
	app.Delete("/api/v1/users/:id", users.DeleteUser)
	app.Post("/api/v1/operation/seedthedb", operation.SeedOperations)
	app.Get("/api/v1/:userid/records", records.GetRecords)
	app.Post("/api/v1/:userid/:operationid/records", records.CreateRecord)
	app.Delete("/api/v1/:userid/records/:id", records.DeleteRecord)
}

func initDatabase() {
	dsn := "host=localhost user=postgres password=password dbname=loanpro port=5432 sslmode=disable"
	var err error
	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connection successfully opened")
	database.DBConn.AutoMigrate(&models.User{}, &models.Record{}, &models.Operation{})
	fmt.Print("Database Migrated")
}

func main() {
	app := fiber.New()
	initDatabase()

	setupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
