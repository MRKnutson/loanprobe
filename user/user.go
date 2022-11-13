package users

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt"

	"github.com/golang-jwt/jwt/v4"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	if err != nil {
		return "error"
	}
	return string(hash)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := claims["user_id"]

	if uid != float64(n) {
		return false
	}

	return true
}

func Validate(c *fiber.Ctx) error {

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Im logged in",
	})
}

func LoginUser(c *fiber.Ctx) error {
	loginAttempt := new(User)

	var err error
	if err = c.BodyParser(loginAttempt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse json",
		})
	}

	var user models.User
	database.DBConn.Raw(`SELECT * FROM users WHERE email = $1`, loginAttempt.Email).Scan(&user)

	if !checkPasswordHash(loginAttempt.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Credentials",
			"data":    nil,
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t, "user": user})
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

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = newUser.ID
	claims["email"] = newUser.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t, "user": newUser})

}

func ValidateToken(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "invalid token id", "data": "nil"})
	}
	var user models.User
	database.DBConn.Raw("SELECT * FROM users WHERE id = $1",
		id).Scan(&user)
	if user.Status == "inactive" {
		c.Status(fiber.StatusConflict)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token, "user": user})
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "invalid token id", "data": "nil"})
	}

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

func LogoutUser(c *fiber.Ctx) error {
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	var user models.User

	database.DBConn.Raw("SELECT * FROM users WHERE id = $1",
		id).Scan(&user)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "invalid token id", "data": "nil"})
	}

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().AddDate(0, 0, -2)

	return c.JSON(fiber.Map{"status": "success", "message": "Success logout", "data": token, "user": user})
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
