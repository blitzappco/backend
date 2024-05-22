package tickets

import (
	"backend/models"
	"backend/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Routes(app *fiber.App) {
	tickets := app.Group("/tickets")

	purchase(tickets)

	tickets.Get("/types", func(c *fiber.Ctx) error {
		city := c.Query("city")

		ticketTypes, err := models.GetTicketTypes(city)
		if err != nil {
			return utils.MessageError(c, "Nu s-au putut gasi tipurile de bilete")
		}

		return c.JSON(ticketTypes)
	})

	tickets.Get("/ticket", func(c *fiber.Ctx) error {
		ticket, err := models.GetTicket(
			fmt.Sprintf("%v", c.Query("ticketID")),
		)

		if err != nil {
			return utils.MessageError(c, err.Error())
		}

		return c.JSON(ticket)
	})

	tickets.Get("/", models.AccountMiddleware, func(c *fiber.Ctx) error {
		accountID := fmt.Sprintf("%v", c.Locals("id"))

		tickets, err := models.GetTickets(accountID)
		if err != nil {
			return utils.MessageError(c, "Nu s-au putut gasi biletele")
		}

		return c.JSON(tickets)
	})

	tickets.Get("/last", models.AccountMiddleware, func(c *fiber.Ctx) error {
		accountID := fmt.Sprintf("%v", c.Locals("id"))
		city := c.Query("city")

		ticket, err := models.GetLastTicket(
			accountID,
			city,
		)

		if err != nil {
			return utils.MessageError(c, "Nu s-a putut gasi biletul")
		}

		show := models.CheckValidity(ticket)

		return c.JSON(bson.M{
			"ticket": ticket,
			"show":   show,
		})
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
