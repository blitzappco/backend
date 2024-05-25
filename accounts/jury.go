package accounts

import (
	"backend/db"
	"backend/models"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func HandleJury(c *fiber.Ctx) error {
	var account models.Account

	account, err := models.GetAccount(bson.M{
		"phone": "+40712345678",
	})

	if err != nil {
		return utils.MessageError(c, err.Error())
	}

	code := "0000"

	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), 10)
	if err != nil {
		return utils.MessageError(c, "A aparut o problema tehnica, incercati mai tarziu.")
	}
	db.Set("code:+40712345678", string(hashedCode))

	token := account.GenAccountToken()

	return c.JSON(bson.M{
		"phone":     "+40712345678",
		"token":     token,
		"newClient": false,
	})
}
