package tickets

import (
	"backend/models"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Routes(app fiber.App) {
	tickets := app.Group("/tickets")

	tickets.Get("/types", func(c *fiber.Ctx) error {
		ticketTypes, err := models.GetTicketTypes("bucuresti")
		if err != nil {
			return utils.MessageError(c, "Nu s-au putut gasi tipurile de bilete")
		}

		return c.JSON(ticketTypes)
	})
}
