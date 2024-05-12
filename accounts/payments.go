package accounts

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/setupintent"

	"github.com/stripe/stripe-go/v78"
	"go.mongodb.org/mongo-driver/bson"
)

func payments(acc fiber.Router) {
	payments := acc.Group("/payments")
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
		var account models.Account
		utils.GetLocals(c, "account", &account)

		fmt.Println(account)

		var body map[string]string
		json.Unmarshal(c.Body(), &body)

		amountToCharge, _ := strconv.Atoi(body["amount"])

		params := &stripe.PaymentIntentParams{
			Customer:      stripe.String(account.StripeCustomerID),
			PaymentMethod: stripe.String(body["paymentMethod"]),
			Amount:        stripe.Int64(int64(amountToCharge)),
			Currency:      stripe.String(string(stripe.CurrencyRON)),
		}
		intent, _ := paymentintent.New(params)

		return c.JSON(bson.M{
			"clientSecret": intent.ClientSecret,
		})
	})
}
