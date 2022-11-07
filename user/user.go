package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelrknutson/loanprogo/database"
	"github.com/michaelrknutson/loanprogo/models"
	"golang.org/x/crypto/bcrypt"

	"time"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(hash)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginUser(c *fiber.Ctx) error {
	loginAttempt := new(User)

	var err error
	if err = c.BodyParser(loginAttempt); err != nil {
		c.SendStatus(503)
	}
	fmt.Println(loginAttempt.Password)

	var user models.User
	database.DBConn.Raw(`SELECT * FROM users WHERE email = $1`, loginAttempt.Email).Scan(&user)

	if checkPasswordHash(loginAttempt.Password, user.PasswordHash) {
		return c.Status(fiber.StatusOK).JSON(user)
	} else {
		return c.Status(fiber.StatusForbidden).JSON("Unauthorized")
	}
}

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DBConn.Raw("SELECT * FROM users WHERE status = 'active'").Scan(&users)
	return c.Status(fiber.StatusOK).JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	database.DBConn.Raw("SELECT * FROM users WHERE id = $1 AND status = 'active'",
		id).Scan(&user)
	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	userInput := new(User)
	newUser := new(models.User)
	var err error
	if err = c.BodyParser(userInput); err != nil {
		c.SendStatus(503)
	}
	newUser.PasswordHash = hashPassword(userInput.Password)
	newUser.Email = userInput.Email

	if err = database.DBConn.Create(&newUser).Error; err != nil {
		c.Status(400).Send([]byte(err.Error()))
	}
	return c.JSON(newUser)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	updatedUser := new(User)
	var err error
	if err = c.BodyParser(updatedUser); err != nil {
		c.SendStatus(503)
	}
	var user models.User
	database.DBConn.Raw("SELECT * FROM users WHERE id = $1",
		id).Scan(&user)
	if user.Status == "inactive" {
		c.Status(400)
	}
	user.Email = updatedUser.Email
	user.PasswordHash = hashPassword(updatedUser.Password)
	if err = database.DBConn.Save(&user).Error; err != nil {
		c.Status(400).Send([]byte(err.Error()))
	}
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var err error

	var user models.User
	database.DBConn.Raw("SELECT * FROM users WHERE id = $1",
		id).Scan(&user)

	t := time.Time(time.Now())
	user.Status = "inactive"
	user.DeletedTime = t
	if err = database.DBConn.Save(&user).Error; err != nil {
		c.Status(400).Send([]byte(err.Error()))
	}

	return c.SendString("User Deleted")
}
