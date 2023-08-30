package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go_linkaja/account"
	"go_linkaja/database"
	"log"
)

func setUpRoutes(app *fiber.App) {
	Account := app.Group("/account")

	Account.Get("/:accountNumber", account.GetSaldo)
	Account.Post("/:fromAccountNumber/transfer", account.PostTransfer)
	// Add your protected endpoints here
}

func main() {
	database.ConnectDb()
	app := fiber.New(fiber.Config{})

	// Add recover middleware
	app.Use(recover.New())

	// Configuring CORS middleware
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "*",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	setUpRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":3000"))
}
