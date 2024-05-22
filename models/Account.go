package models

import (
	"backend/db"
	"backend/env"
	"backend/utils"
	"strings"
	"time"

	sj "github.com/brianvoe/sjwt"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/customer"
	"go.mongodb.org/mongo-driver/bson"
)

type Account struct {
	ID               string `bson:"id" json:"id"`
	Phone            string `bson:"phone" json:"phone"`
	FirstName        string `bson:"firstName" json:"firstName"`
	LastName         string `bson:"lastName" json:"lastName"`
	StripeCustomerID string `bson:"stripeCustomerID" json:"stripeCustomerID"`

	PaymentMethods []PaymentMethod `bson:"paymentMethods" json:"paymentMethods"`
	Trips          []Trip          `bson:"trips" json:"trips"`
}

type PaymentMethod struct {
	ID    string `bson:"id" json:"id"`
	Type  string `bson:"type" json:"type"`
	Icon  string `bson:"icon" json:"icon"`
	Title string `bson:"title" json:"title"`
}

type Trip struct {
	PlaceID       string `bson:"placeID" json:"placeID"`
	MainText      string `bson:"mainText" json:"mainText"`
	SecondaryText string `bson:"secondaryText" json:"secondaryText"`
	Type          string `bson:"type" json:"type"`
}

func (account Account) GenAccountToken() string {
	claims, _ := sj.ToClaims(account)
	claims.SetExpiresAt(time.Now().Add(365 * 24 * time.Hour))
	// 1 year = 365 days * 24 hours in a day

	token := claims.Generate(env.JWTKey)
	return token
}

func ParseAccountToken(token string) (Account, error) {
	hasVerified := sj.Verify(token, env.JWTKey)

	if !hasVerified {
		return Account{}, nil
	}

	claims, _ := sj.Parse(token)
	err := claims.Validate()
	account := Account{}
	claims.ToStruct(&account)

	return account, err
}

func AccountMiddleware(c *fiber.Ctx) error {
	var token string

	authHeader := c.Get("Authorization")

	if string(authHeader) != "" && strings.HasPrefix(string(authHeader), "Bearer") {

		tokens := strings.Fields(string(authHeader))
		if len(tokens) == 2 {
			token = tokens[1]
		}
		if token == "" {
			return utils.MessageError(c, "no token")
		}

		account, err := ParseAccountToken(token)
		if err != nil {
			return utils.MessageError(c, "A aparut o eroare")
		}

		c.Locals("id", account.ID)
		utils.SetLocals(c, "account", account)
	}

	if token == "" {
		return utils.MessageError(c, "no token")
	}

	return c.Next()
}

func (account *Account) Create() error {
	// generating ID
	account.ID = utils.GenID(6)

	// creating account
	_, err := db.Accounts.InsertOne(db.Ctx, account)

	return err
}

func (account *Account) CreateStripeCustomer() error {
	fullName := account.LastName + " " + account.FirstName

	// creating a stripe account
	params := &stripe.CustomerParams{
		Phone: &account.Phone,
		Name:  &fullName,
	}
	c, err := customer.New(params)

	if err != nil {
		return err
	}

	account.StripeCustomerID = c.ID
	return nil
}

func (account *Account) AddPaymentMethod(pm PaymentMethod) error {
	tempAccount, err := GetAccount(bson.M{"id": account.ID})

	if err != nil {
		return err
	}

	account.PaymentMethods = append(tempAccount.PaymentMethods, pm)

	err = UpdateAccount(account.ID, bson.M{
		"paymentMethods": account.PaymentMethods,
	})

	if err != nil {
		return err
	}

	return nil
}

func UpdateAccount(id string, updates interface{}) error {
	_, err := db.Accounts.UpdateOne(
		db.Ctx,
		bson.M{"id": id},
		bson.M{
			"$set": updates,
		},
	)

	return err
}

func GetAccount(query interface{}) (Account, error) {
	var account Account

	err := db.Accounts.FindOne(
		db.Ctx,
		query,
	).Decode(&account)

	return account, err
}

func CheckAccount(phone string) (bool, Account) {
	var account Account

	err := db.Accounts.FindOne(
		db.Ctx, bson.M{
			"phone": phone,
		},
	).Decode(&account)

	if err != nil {
		return false, Account{}
	} else {
		return true, account
	}
}

func AddTrip(accountID string, trip Trip) (Account, error) {
	account, err := GetAccount(bson.M{
		"id": accountID,
	})

	if err != nil {
		return account, err
	}

	newTrips := append(account.Trips, trip)
	if len(newTrips) > 10 {
		newTrips = newTrips[1:]
	}

	err = UpdateAccount(accountID, bson.M{
		"trips": newTrips,
	})

	if err != nil {
		return account, err
	}

	account.Trips = newTrips
	return account, err
}
