package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/michaelrknutson/loanprogo/database"
	"github.com/michaelrknutson/loanprogo/initializers"
	"github.com/michaelrknutson/loanprogo/middleware"
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
	app.Get("/api/v1/users/validate", users.Validate)
	app.Post("/api/v1/users/", users.CreateUser)
	app.Post("/api/v1/users/login", users.LoginUser)
	app.Post("/api/v1/users/:id/logout", middleware.RequireAuth(), users.LogoutUser)
	app.Put("/api/v1/users/:id", middleware.RequireAuth(), users.UpdateUser)
	app.Delete("/api/v1/users/:id", users.DeleteUser)
	app.Post("/api/v1/operation/seedthedb", operation.SeedOperations)
	app.Get("/api/v1/:userid/records", middleware.RequireAuth(), records.GetRecords)
	app.Get("/api/v1/validate", users.Validate)
	app.Get("/api/v1/:id/auth/validate_token", middleware.RequireAuth(), users.ValidateToken)
	app.Post("/api/v1/:userid/:operationid/records", middleware.RequireAuth(), records.CreateRecord)
	app.Delete("/api/v1/:userid/records/:id", middleware.RequireAuth(), records.DeleteRecord)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
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

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := fiber.New()
	initDatabase()

	app.Use(cors.New())
	// app.Use(middleware.RequireAuth)

	setupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}
