package accounts

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func trips(acc fiber.Router) {
	trips := acc.Group("/trips")

	trips.Get("/", models.AccountMiddleware, func(c *fiber.Ctx) error {
		accountID := fmt.Sprintf("%v", c.Locals("id"))

		account, err := models.GetAccount(bson.M{
			"id": accountID,
		})

		if err != nil {
			return utils.MessageError(c, err.Error())
		}

		return c.JSON(
			account.Trips,
		)
	})

	trips.Put("/", models.AccountMiddleware, func(c *fiber.Ctx) error {
		var trip models.Trip
		json.Unmarshal(c.Body(), &trip)

		accountID := fmt.Sprintf("%v", c.Locals("id"))

		account, err := models.AddTrip(accountID, trip)

		if err != nil {
			return utils.MessageError(c, err.Error())
		}

		return c.JSON(account.Trips)
	})
}
