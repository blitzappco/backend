package accounts

import "github.com/gofiber/fiber/v2"

func Routes(app *fiber.App) {
	acc := app.Group("/accounts")

	onboarding(acc)
	payments(acc)
	trips(acc)
}
