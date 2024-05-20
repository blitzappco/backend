package accounts

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/setupintent"

	"github.com/stripe/stripe-go/v78"
	"go.mongodb.org/mongo-driver/bson"
)

func payments(acc fiber.Router) {
	payments := acc.Group("/payments")

	payments.Get("/methods", models.AccountMiddleware, func(c *fiber.Ctx) error {
		accountID := fmt.Sprintf("%v", c.Locals("id"))

		account, err := models.GetAccount(
			bson.M{
				"id": accountID,
			},
		)

		if err != nil {
			return utils.MessageError(c, err.Error())
		}

		return c.JSON(account.PaymentMethods)
	})

	payments.Post("/setup-intent", models.AccountMiddleware, func(c *fiber.Ctx) error {
		var account models.Account
		utils.GetLocals(c, "account", &account)

		params := &stripe.SetupIntentParams{
			AutomaticPaymentMethods: &stripe.SetupIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
			Customer: stripe.String(account.StripeCustomerID),
		}
		intent, _ := setupintent.New(params)

		return c.JSON(bson.M{
			"clientSecret": intent.ClientSecret,
		})
	})

	payments.Post("/setup-confirm", models.AccountMiddleware, func(c *fiber.Ctx) error {
		var account models.Account
		utils.GetLocals(c, "account", &account)

		var body map[string]string
		json.Unmarshal(c.Body(), &body)

		params := &stripe.CustomerRetrievePaymentMethodParams{
			Customer: stripe.String(account.StripeCustomerID),
		}
		result, err := customer.RetrievePaymentMethod(body["paymentMethod"], params)

		if err != nil {
			return utils.MessageError(c, "A aparut o eroare")
		}

		// figuring out what type of payment method it is
		switch result.Type {
		case "card":
			err = account.AddPaymentMethod(models.PaymentMethod{
				ID:    body["paymentMethod"],
				Type:  "card",
				Icon:  result.Card.DisplayBrand,
				Title: result.Card.Last4,
			})

		}
		if err != nil {
			return utils.MessageError(c, "Nu a mers")
		}

		token := account.GenAccountToken()

		return c.JSON(bson.M{
			"token": token,
		})
	})

	payments.Post("/payment-intent", models.AccountMiddleware, func(c *fiber.Ctx) error {
		accountID := fmt.Sprintf("%v", c.Locals("id"))

		account, err := models.GetAccount(bson.M{"id": accountID})
		if err != nil {
			return utils.MessageError(c, "Nu s-a putut gasi contul")
		}

		var body struct {
			Amount        int `json:"amount"`
			PaymentMethod int `json:"paymentMethod"`
		}
		json.Unmarshal(c.Body(), &body)

		params := &stripe.PaymentIntentParams{
			Customer: stripe.String(account.StripeCustomerID),
			PaymentMethod: stripe.String(
				account.PaymentMethods[body.PaymentMethod].ID,
			),
			Amount:   stripe.Int64(int64(body.Amount)),
			Currency: stripe.String(string(stripe.CurrencyRON)),
		}
		intent, _ := paymentintent.New(params)

		return c.JSON(bson.M{
			"clientSecret":  intent.ClientSecret,
			"paymentIntent": intent.ID,
		})
	})
}
