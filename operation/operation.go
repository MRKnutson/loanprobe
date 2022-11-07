package operation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/michaelrknutson/loanprogo/database"
	"github.com/michaelrknutson/loanprogo/models"
)

func SeedOperations(c *fiber.Ctx) error {
	var operations = []models.Operation{{Type: "Addition", Cost: 1.25}, {Type: "Subtraction", Cost: 2.25}, {Type: "Multiplication", Cost: 5.00}, {Type: "Division", Cost: 6.25}, {Type: "Square_root", Cost: 7.00}, {Type: "Random_String", Cost: 15.00}}
	database.DBConn.Create(&operations)

	return c.JSON(operations)
}
