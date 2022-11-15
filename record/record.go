package records

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/michaelrknutson/loanprogo/database"
	"github.com/michaelrknutson/loanprogo/models"
	"github.com/nouney/randomstring"
)

type Record struct {
	FirstNumber       float64          `json:"firstnumber"`
	SecondNumber      float64          `json:"secondnumber"`
	Amount            float64          `json:"amount"`
	UserBalance       float64          `json:"userbalance"`
	OperationResponse string           `json:"operationresponse"`
	Operation         models.Operation `json:"operation"`
	User              models.User      `json:"user"`
}

func checkBalance(lastRecord models.Record) (balance float64) {
	if lastRecord.ID > 0 {
		balance = lastRecord.UserBalance
	} else {
		balance = 300.00
	}
	return
}

func checkOperation(operation models.Operation, balance float64, firstNumber float64, secondNumber float64) (response string, newBalance float64) {
	if operation.Cost > balance {
		response = "no-op"
		newBalance = balance
	} else {
		response = performOperation(operation.Type, firstNumber, secondNumber)
		newBalance = balance - operation.Cost
	}
	return
}

func generateStrings() (randomString string) {
	response, err := http.Get("https://www.random.org/strings/?num=1&len=16&digits=on&upperalpha=on&loweralpha=on&unique=on&format=plain&rnd=new")
	if err != nil {
		randomString = randomstring.Generate(12)
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal((err))
	}
	randomString = fmt.Sprintf("%x", responseData)
	return
}

func performOperation(operation string, firstNumber float64, secondNumber float64) (operationResponse string) {
	switch operation {
	case "Addition":
		operationResponse = fmt.Sprint(firstNumber + secondNumber)
	case "Subtraction":
		operationResponse = fmt.Sprint(firstNumber - secondNumber)
	case "Multiplication":
		operationResponse = fmt.Sprint(firstNumber * secondNumber)
	case "Division":
		operationResponse = fmt.Sprint(firstNumber / secondNumber)
	case "Square_root":
		operationResponse = fmt.Sprint(math.Sqrt(firstNumber))
	case "Random_String":
		operationResponse = generateStrings()
	}
	return
}

func GetRecords(c *fiber.Ctx) error {

	userId := c.Params("userid")
	var records []models.Record
	database.DBConn.Raw("SELECT * FROM records WHERE deleted_time IS NULL AND user_refer = $1",
		userId).Scan(&records)
	return c.Status(fiber.StatusOK).JSON(records)
}

func CreateRecord(c *fiber.Ctx) error {
	newRecord := new(Record)
	if err := c.BodyParser(newRecord); err != nil {
		c.SendStatus(200)
	}
	userId := c.Params("userid")
	operationId := c.Params("operationid")

	var user models.User
	database.DBConn.Raw("SELECT * FROM users WHERE id = $1 AND status = 'active'",
		userId).Scan(&user)
	var operation models.Operation
	database.DBConn.Raw("SELECT * FROM operations WHERE id = $1",
		operationId).Scan(&operation)
	var lastRecord models.Record
	database.DBConn.Raw("SELECT * FROM records WHERE user_refer = $1 ORDER BY id DESC LIMIT 1;",
		userId).Scan(&lastRecord)

	currentBalance := checkBalance(lastRecord)

	response, newBalance := checkOperation(operation, currentBalance, newRecord.FirstNumber, newRecord.SecondNumber)

	fullRecord := new(models.Record)
	fullRecord.UserRefer = int(user.ID)
	fullRecord.OperationRefer = int(operation.ID)
	fullRecord.Amount = operation.Cost
	fullRecord.UserBalance = newBalance
	fullRecord.OperationResponse = response
	fullRecord.Operation = operation
	fullRecord.User = user

	var err error
	if err = database.DBConn.Create(fullRecord).Error; err != nil {
		c.Status(400).Send([]byte(err.Error()))
	}

	return c.JSON(fullRecord)
}

func DeleteRecord(c *fiber.Ctx) error {
	userId := c.Params("userid")
	recordId := c.Params("id")
	// database.DBConn.Raw("UPDATE users SET status = 'inactive', DeletedAt = $1 WHERE id = 2")
	var err error

	var record models.Record
	database.DBConn.Raw("SELECT * FROM records WHERE user_refer = $1 AND id = $2",
		userId, recordId).Scan(&record)

	t := time.Now()
	record.DeletedTime = t
	if err = database.DBConn.Save(&record).Error; err != nil {
		c.Status(400).Send([]byte(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(record)
}
