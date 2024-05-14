package models

import (
	"backend/db"
	"backend/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ticket struct {
	ID        string `bson:"id" json:"id"`
	AccountID string `bson:"accountID" json:"accountID"`
	City      string `bson:"city" json:"city"`

	Mode   string  `bson:"mode" json:"mode"`
	Fare   float64 `bson:"fare" json:"fare"`
	Trips  int     `bson:"trips" json:"trips"`
	Expiry string  `bson:"expiry" json:"expiry"`

	ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func GetTickets(accountID string) ([]Ticket, error) {
	cursor, err := db.Tickets.Find(db.Ctx, bson.M{
		"accountID": accountID,
	})

	if err != nil {
		return []Ticket{}, err
	}

	tickets := []Ticket{}

	err = cursor.All(db.Ctx, &tickets)
	if len(tickets) == 0 {
		tickets = []Ticket{}
	}

	return tickets, err
}

func GetTicket(ticketID string) (Ticket, error) {
	ticket := Ticket{}

	err := db.Tickets.
		FindOne(db.Ctx, bson.M{"id": ticketID}).
		Decode(&ticket)

	return ticket, err
}

func ChangeTicket(ticketID string, updates interface{}) (Ticket, error) {
	ticket := Ticket{}

	err := db.Tickets.FindOneAndUpdate(
		db.Ctx,
		bson.M{"id": ticketID},
		bson.M{
			"$set": updates,
		},
	).Decode(&ticket)

	return ticket, err
}

func ConfirmTicket(ticketID string, accountID string) (Ticket, error) {
	ticket, err := ChangeTicket(ticketID,
		bson.M{
			"accountID": accountID,
		},
	)

	ticket.AccountID = accountID

	return ticket, err
}

func GetLastTicket(accountID string) (Ticket, error) {
	tickets := []Ticket{}

	cursor, err := db.Tickets.Find(db.Ctx,
		bson.M{"accountID": accountID},
		options.Find().SetLimit(1).
			SetSort(bson.M{"createdAt": -1}))

	if err != nil {
		return tickets[0], err
	}

	err = cursor.All(db.Ctx, &tickets)
	if err != nil {
		return tickets[0], err
	}

	return tickets[0], nil
}

func (ticket *Ticket) Create(tt TicketType) error {
	ticket.ID = utils.GenID(12)

	loc, _ := time.LoadLocation("Europe/Bucharest")
	ticket.CreatedAt = time.Now().In(loc)

	ticket.City = tt.City
	ticket.Mode = tt.Mode
	ticket.Fare = tt.Fare
	ticket.Trips = tt.Trips
	ticket.Expiry = tt.Expiry

	_, err := db.Tickets.InsertOne(db.Ctx, ticket)
	return err
}
