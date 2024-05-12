package main

import (
	"backend/accounts"
	"backend/db"
	"backend/env"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stripe/stripe-go/v78"
)

func main() {
	app := fiber.New()

	stripe.Key = env.StripeKey

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	db.InitDB()
	db.InitCache()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("PONG")
	})

	accounts.Routes(app)
	app.Listen(":6969")
}
