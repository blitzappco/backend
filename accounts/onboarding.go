package accounts

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func onboarding(acc fiber.Router) {
	onboarding := acc.Group("/onboarding")

	onboarding.Post("/test", func(c *fiber.Ctx) error {
		account, _ := models.GetAccount(bson.M{"phone": "+40723010405"})

		return c.JSON(account)
	})

	onboarding.Post("", func(c *fiber.Ctx) error {
		var body map[string]string
		json.Unmarshal(c.Body(), &body)
		phone := body["phone"]

		if phone == "+40712345678" {
			return HandleJury(c)
		}

		exists, account := models.CheckAccount(phone)
		fmt.Println(account)

		// generating code
		code := utils.GenCode(4)
		fmt.Println(code)

		// sending verification code on sms
		// err := utils.SendSMS(phone, code)
		// ok := true
		// if err != nil {
		//   return utils.MessageError(c, fmt.Sprintf("Codul nu a putut fi trimis.: %v", err))
		// }

		// setting verification code in redis
		hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), 10)
		if err != nil {
			return utils.MessageError(c, "A aparut o problema tehnica, incercati mai tarziu.")
		}
		db.Set("code:"+phone, string(hashedCode))

		if !exists {
			account.Create()
		}

		token := account.GenAccountToken()

		return c.JSON(bson.M{
			"phone":     phone,
			"token":     token,
			"newClient": !exists,
		})
	})

	onboarding.Post("/verify-code", models.AccountMiddleware, func(c *fiber.Ctx) error {
		var body map[string]string
		json.Unmarshal(c.Body(), &body)
		phone := body["phone"]
		code := body["code"]

		// getting code from redis
		hashedCode, _ := db.Get("code:" + phone)

		if bcrypt.CompareHashAndPassword(
			[]byte(hashedCode), []byte(code)) != nil {
			return utils.MessageError(c, "Codul introdus nu este corect")
		}

		// get account from token
		account := models.Account{}
		utils.GetLocals(c, "account", &account)

		// update account with phone
		err := models.UpdateAccount(
			account.ID,
			bson.M{
				"phone": phone,
			},
		)
		if err != nil {
			return utils.MessageError(c, "Eroare")
		}

		// update local account with phone
		account.Phone = phone

		// gen token
		token := account.GenAccountToken()

		return c.JSON(bson.M{
			"account": account,
			"token":   token,
		})
	})

	onboarding.Post("/name", models.AccountMiddleware, func(c *fiber.Ctx) error {
		var body map[string]string
		json.Unmarshal(c.Body(), &body)
		firstName := body["firstName"]
		lastName := body["lastName"]

		account := models.Account{}
		utils.GetLocals(c, "account", &account)

		// update local account with name
		account.FirstName = firstName
		account.LastName = lastName

		err := account.CreateStripeCustomer()
		if err != nil {
			return utils.MessageError(c, "Nu s-a putut verifica")
		}

		// update account with name
		err = models.UpdateAccount(
			account.ID,
			bson.M{
				"firstName":        firstName,
				"lastName":         lastName,
				"stripeCustomerID": account.StripeCustomerID,
			},
		)
		if err != nil {
			return utils.MessageError(c, "Eroare")
		}

		// generating the account token
		token := account.GenAccountToken()

		return c.JSON(bson.M{
			"account": account,
			"token":   token,
		})
	})
}
