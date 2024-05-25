package models

import (
	"backend/db"

	"go.mongodb.org/mongo-driver/bson"
)

type TicketType struct {
	ID   string `bson:"id" json:"id"`
	City string `bson:"city" json:"city"`

	Mode   string `bson:"mode" json:"mode"`
	Fare   int    `bson:"fare" json:"fare"`
	Trips  int    `bson:"trips" json:"trips"`
	Expiry string `bson:"expiry" json:"expiry"`
}

func GetTicketTypes(city string) ([]TicketType, error) {
	cursor, err := db.TicketTypes.Find(db.Ctx, bson.M{
		"city": city,
	})

	if err != nil {
		return []TicketType{}, err
	}

	ticketTypes := []TicketType{}

	err = cursor.All(db.Ctx, &ticketTypes)
	if len(ticketTypes) == 0 {
		ticketTypes = []TicketType{}
	}

	return ticketTypes, err
}

func GetTicketType(id string) (TicketType, error) {
	var ticketType TicketType

	err := db.TicketTypes.FindOne(
		db.Ctx,
		bson.M{
			"id": id,
		},
	).Decode(&ticketType)

	return ticketType, err
}
