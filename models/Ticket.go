package models

import (
	"backend/db"
	"backend/utils"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ticket struct {
	ID        string `bson:"id" json:"id"`
	AccountID string `bson:"accountID" json:"accountID"`
	City      string `bson:"city" json:"city"`
	Name      string `bson:"name" json:"name"`

	Mode   string `bson:"mode" json:"mode"`
	Fare   int    `bson:"fare" json:"fare"`
	Trips  int    `bson:"trips" json:"trips"`
	Expiry string `bson:"expiry" json:"expiry"`

	PaymentIntent string `bson:"paymentIntent" json:"paymentIntent"`
	Confirmed     bool   `bson:"confirmed" json:"confirmed"`

	ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

func GetTickets(accountID string) ([]Ticket, error) {
	cursor, err := db.Tickets.Find(db.Ctx, bson.M{
		"accountID": accountID,
		"confirmed": true,
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

func GetTicketPaymentIntent(paymentIntent string) (Ticket, error) {
	ticket := Ticket{}

	err := db.Tickets.FindOne(db.Ctx, bson.M{"paymentIntent": paymentIntent}).Decode(&ticket)

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

func ConfirmTicket(paymentIntent string) error {
	_, err := db.Tickets.
		UpdateOne(
			db.Ctx,
			bson.M{"paymentIntent": paymentIntent},
			bson.M{
				"$set": bson.M{
					"confirmed": true,
				},
			},
		)

	return err
}

func GetLastTicket(accountID string, city string) (Ticket, error) {
	tickets := []Ticket{}

	cursor, err := db.Tickets.Find(db.Ctx,
		bson.M{"accountID": accountID, "city": city},
		options.Find().SetLimit(1).
			SetSort(bson.M{"createdAt": -1}))

	if err != nil {
		return Ticket{}, err
	}

	err = cursor.All(db.Ctx, &tickets)
	if err != nil {
		return Ticket{}, err
	}

	if len(tickets) != 0 {
		return tickets[0], nil
	} else {
		return Ticket{Fare: 0.0}, nil
	}
}

func (ticket *Ticket) Create(tt TicketType, accountID string, name string) error {
	ticket.ID = utils.GenID(12)

	now := time.Now().UTC()
	ticket.CreatedAt = now

	ticket.Name = name
	ticket.AccountID = accountID
	ticket.City = tt.City
	ticket.Mode = tt.Mode
	ticket.Fare = tt.Fare
	ticket.Trips = tt.Trips
	ticket.Expiry = tt.Expiry

	_, err := db.Tickets.InsertOne(db.Ctx, ticket)
	return err
}

func Validate(ticketID string) (Ticket, bool, error) {
	ticket, err := GetTicket(ticketID)
	if err != nil {
		return ticket, false, err
	}

	switch ticket.City {
	case "bucuresti":
		return bucharestValidate(ticket)
	case "ploiesti":
		return ploiestiValidate(ticket)
	default:
		return ticket, false, errors.New("acest bilet sau abonament nu poate fi inregistrat in sistem")
	}

}

func ploiestiValidate(ticket Ticket) (Ticket, bool, error) {
	return ticket, false, nil
}

func bucharestValidate(ticket Ticket) (Ticket, bool, error) {
	var err error
	if ticket.Trips < 0 { // it is a pass
		if ticket.ExpiresAt.IsZero() {
			expiresAt := utils.ConvertExpiry(ticket.Expiry, time.Now().UTC())
			ticket, err := ChangeTicket(
				ticket.ID,
				bson.M{
					"expiresAt": expiresAt,
				},
			)
			ticket.ExpiresAt = expiresAt

			return ticket, true, err
		} else {
			if ticket.ExpiresAt.After(time.Now()) {
				return ticket, true, err
			} else {
				return ticket, false, err
			}
		}
	} else { // it is a ticket
		if ticket.ExpiresAt.After(time.Now()) {
			return ticket, true, err
		} else {
			if ticket.Trips > 0 {
				expiresAt := utils.ConvertExpiry(ticket.Expiry, time.Now())
				trips := ticket.Trips - 1
				ticket, err := ChangeTicket(
					ticket.ID,
					bson.M{
						"expiresAt": expiresAt,
						"trips":     trips,
					},
				)

				ticket.ExpiresAt = expiresAt
				ticket.Trips = trips

				return ticket, true, err
			} else {
				return ticket, false, err
			}
		}
	}
}

func CheckValidity(ticket Ticket) bool {
	switch ticket.City {
	case "bucuresti":
		return bucharestCheck(ticket)
	case "ploiesti":
		return ploiestiCheck(ticket)
	default:
		return false
	}
}

func bucharestCheck(ticket Ticket) bool {
	show := false

	// see if it should actually show it
	if ticket.Trips < 0 { // it's a pass, show it if it is still available
		if ticket.ExpiresAt.IsZero() {
			show = true
		} else if ticket.ExpiresAt.After(time.Now().UTC()) {
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

	return show
}

func ploiestiCheck(_ Ticket) bool {
	return false
}
