package tickets

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Routes(app *fiber.App) {
	tickets := app.Group("/tickets")

	purchase(tickets)

	tickets.Get("/types", func(c *fiber.Ctx) error {
		ticketTypes, err := models.GetTicketTypes("bucuresti")
		if err != nil {
			return utils.MessageError(c, "Nu s-au putut gasi tipurile de bilete")
		}

		return c.JSON(ticketTypes)
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

		ticket, err := models.GetLastTicket(
			accountID,
		)

		if err != nil {
			return utils.MessageError(c, "Nu s-a putut gasi biletul")
		}

		show := false

		// see if it should actually show it
		if ticket.Trips < 0 { // it's a pass, show it if it is still available
			if ticket.ExpiresAt.After(time.Now().UTC()) {
				show = true
			}
		} else { // it's a ticket, we'll have to see the number of trips left
			if ticket.Trips > 0 {
				show = true
			} else {
				if ticket.ExpiresAt.After(time.Now().UTC()) {
					show = true
				}
			}
		}

		if !ticket.Confirmed {
			show = false
		}

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
