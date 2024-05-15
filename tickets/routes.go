package tickets

import (
	"backend/models"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Routes(app *fiber.App) {
	tickets := app.Group("/tickets")

	tickets.Get("/types", func(c *fiber.Ctx) error {
		ticketTypes, err := models.GetTicketTypes("bucuresti")
		if err != nil {
			return utils.MessageError(c, "Nu s-au putut gasi tipurile de bilete")
		}

		return c.JSON(ticketTypes)
	})

	tickets.Post("/purchase-intent", models.AccountMiddleware, func(c *fiber.Ctx) error {
		typeID := c.Query("typeID")
		ticketType, err := models.GetTicketType(typeID)
		if err != nil {
			return utils.MessageError(c, "Nu s-a putut gasi tipul de bilet")
		}

		var ticket models.Ticket
		err = ticket.Create(ticketType)
		if err != nil {
			return utils.MessageError(c, "Nu s-a putut crea biletul")
		}

		return c.JSON(
			bson.M{
				"fare":     ticketType.Fare,
				"ticketID": ticket.ID,
			},
		)
	})

	tickets.Post("/purchase-confirm", models.AccountMiddleware, func(c *fiber.Ctx) error {
		var account models.Account
		utils.GetLocals(c, "account", &account)

		ticketID := c.Query("ticketID")

		ticket, err := models.ConfirmTicket(ticketID, account.ID)
		if err != nil {
			return utils.MessageError(c, "Nu s-a putut gasi biletul")
		}

		return c.JSON(ticket)
	})

	tickets.Post("/validate", models.AccountMiddleware, func(c *fiber.Ctx) error {
		ticketID := c.Query("ticketID")

		var ticket models.Ticket
		ticket, valid, err := models.Validate(ticketID)

		if err != nil {
			return utils.MessageError(c, err.Error())
		}

		return c.JSON(
			bson.M{
				"ticket": ticket,
				"valid":  valid,
			},
		)
	})
}
