package tickets

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v78"
	"go.mongodb.org/mongo-driver/bson"
)

func purchase(tickets fiber.Router) {
	tickets.Post("/purchase-intent", models.AccountMiddleware, func(c *fiber.Ctx) error {

		var body map[string]string
		json.Unmarshal(c.Body(), &body)

		typeID := c.Query("typeID")
		ticketType, err := models.GetTicketType(typeID)
		if err != nil {
			return utils.MessageError(c, "Nu s-a putut gasi tipul de bilet")
		}

		var account models.Account
		utils.GetLocals(c, "account", &account)

		var ticket models.Ticket
		err = ticket.Create(ticketType, account.ID, body["name"])
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

	tickets.Post("/purchase-attach", models.AccountMiddleware, func(c *fiber.Ctx) error {
		var body struct {
			PaymentIntent string `json:"paymentIntent"`
			TicketID      string `json:"ticketID"`
		}
		json.Unmarshal(c.Body(), &body)

		ticket, err := models.ChangeTicket(
			body.TicketID,
			bson.M{
				"paymentIntent": body.PaymentIntent,
			},
		)
		ticket.PaymentIntent = body.PaymentIntent

		if err != nil {
			return utils.MessageError(c, err.Error())
		}
		return c.JSON(ticket)
	})

	tickets.Post("/purchase-confirm", func(c *fiber.Ctx) error {
		event := stripe.Event{}
		json.Unmarshal(c.Body(), &event)

		paymentIntent := fmt.Sprintf("%v", event.Data.Object["id"])

		switch event.Type {
		case "payment_intent.succeeded":
			err := models.ConfirmTicket(paymentIntent)
			if err != nil {
				return utils.MessageError(c, err.Error())
			}
		}

		return c.SendString("hello")
	})
}
