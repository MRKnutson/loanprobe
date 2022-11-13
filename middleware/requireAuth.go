package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func RequireAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("SECRET")),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {

	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

// tokenString := c.Cookies("Authorization")

// token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 	}

// 	return []byte(os.Getenv("SECRET")), nil
// })

// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 	if float64(time.Now().Unix()) > claims["exp"].(float64) {
// 		c.SendStatus(fiber.StatusUnauthorized)
// 	}
// 	var user models.User
// 	database.DBConn.First(&user, claims["sub"])

// 	if user.ID == 0 {
// 		c.SendStatus(fiber.StatusBadRequest)
// 	}

// 	c.Set("user", user.Email)
// 	c.Next()
// } else {
// 	c.SendStatus(fiber.StatusUnauthorized)
// }

// // return jwtware.New(jwtware.Config{
// // 	SigningKey: []byte(os.Getenv("SECRET")),
// // })
// fmt.Println("In Middleware")
// return c.Next().Error()
// }
