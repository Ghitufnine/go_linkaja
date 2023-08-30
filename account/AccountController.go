package account

import (
	"github.com/gofiber/fiber/v2"
	"go_linkaja/database"
	"strconv"
)

type SaldoStruct struct { // Structure JSON for GetSaldoQuery function
	AccountNumber string `gorm:"column:account_number" json:"account_number"`
	CustomerName  string `gorm:"column:customer_name" json:"customer_name"`
	Balance       int    `gorm:"column:balance" json:"balance"`
}

func GetSaldoQuery(accountNumber int) []SaldoStruct { // function query get saldo from account_number
	result := []SaldoStruct{}

	database.DBConn.Unscoped().
		Table("account").
		Select("account.account_number AS account_number,"+
			"customer.name AS customer_name,"+
			"account.balance").
		Joins("left join customer on account.customer_number = customer.customer_number").
		Where("account.account_number = ?", accountNumber).
		Find(&result)

	return result
}

func GetSaldo(c *fiber.Ctx) error {
	accountNumber := c.Params("accountNumber") //Params {account_number}
	if accountNumber == "" {                   // Check if accountNumber is not fill
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Data not complete!",
		})
	}

	accountNumberInt, _ := strconv.Atoi(accountNumber) // Change type accountNumber from string into integer
	getSaldo := GetSaldoQuery(accountNumberInt)        // Querying getSaldo using accountNumber
	if len(getSaldo) == 0 {                            // check if data not found or account_number is invalid
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "No account_number found",
		})
	} else {
		getSaldoReturn := getSaldo[0]
		return c.Status(200).JSON(getSaldoReturn)
	}
}

type AmountUpdateForAccount struct {
	Amount int `gorm:"column:balance" json:"amount"`
}

func GetLastBalanceByAccountNumber(accountNumber string) int {
	balance := 0
	type Balance struct {
		AmountBalance int `gorm:"column:balance" json:"amount_balance"`
	}
	amountBalance := Balance{}
	database.DBConn.Unscoped().
		Table("account").
		Select("balance").
		Where("account_number = ?", accountNumber).
		Find(&amountBalance)

	balance = amountBalance.AmountBalance

	return balance
}

func CheckExistAccountNumberQuery(accountNumber string) int {
	count := 0

	database.DBConn.Unscoped().
		Table("account").
		Where("account_number = ?", accountNumber).
		Count(&count)

	return count
}

func PostTransfer(c *fiber.Ctx) error {
	fromAccountNumber := c.Params("fromAccountNumber") //param post for {from_account_number}
	if fromAccountNumber == "" {                       // Check if fromAccountNumber is not fill
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Data not complete or No from_account_number found!",
		})
	}

	type RequestBody struct { // Mapping JSON Request Body
		ToAccountNumber string `json:"to_account_number"`
		Amount          int    `json:"amount"`
	}

	requestBody := new(RequestBody)
	err := c.BodyParser(requestBody)
	if err != nil {
		return err
	}

	checkExistFromAccountNumber := CheckExistAccountNumberQuery(fromAccountNumber)
	if checkExistFromAccountNumber == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Account number: " + fromAccountNumber + " not found",
		})
	}

	// Filling Parse JSON from request body into each object
	toAccountNumber := requestBody.ToAccountNumber
	amount := requestBody.Amount

	checkExistToAccountNumber := CheckExistAccountNumberQuery(toAccountNumber)
	if checkExistToAccountNumber == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Account number: " + toAccountNumber + " not found",
		})
	}

	getFromAccountNumberBalance := GetLastBalanceByAccountNumber(fromAccountNumber)
	lastAmountFromAccountBalance := getFromAccountNumberBalance
	if amount > lastAmountFromAccountBalance {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": "Balance not sufficient",
		})
	}
	getToAccountNumberBalance := GetLastBalanceByAccountNumber(toAccountNumber)
	lasBalanceToAccountNumber := getToAccountNumberBalance
	newAmountToAccountNumber := lasBalanceToAccountNumber + amount
	newAmountFromAccountNumber := lastAmountFromAccountBalance - amount

	// Data New Amount for to_account_number
	DataAmountUpdateToAccountNumber := AmountUpdateForAccount{
		Amount: newAmountToAccountNumber,
	}

	// Data New Amount for from_account_number
	DataAmountUpdateFromAccountNumber := AmountUpdateForAccount{
		Amount: newAmountFromAccountNumber,
	}
	// Transaction Update from Transfer
	transactionTransfer := TransactionTransfer(toAccountNumber, DataAmountUpdateToAccountNumber, fromAccountNumber, DataAmountUpdateFromAccountNumber)
	if transactionTransfer != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  false,
			"message": transactionTransfer.Error(),
		})
	} else {
		return c.Status(201).JSON(fiber.StatusCreated)
	}
}

func TransactionTransfer(toAccountNumber string, DataAmountUpdateToAccountNumber AmountUpdateForAccount, fromAccountNumber string, DataAmountUpdateFromAccountNumber AmountUpdateForAccount) error {
	tx := database.DBConn.Begin()

	// Transaction new balance for {to_account_number}
	if err := tx.
		Table("account").
		Where("account_number = ?", toAccountNumber).
		Updates(&DataAmountUpdateToAccountNumber).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Transaction new balance for {from_account_number}
	if err := tx.
		Table("account").
		Where("account_number = ?", fromAccountNumber).
		Updates(&DataAmountUpdateFromAccountNumber).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
