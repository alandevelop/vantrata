package dto

import "github.com/volatiletech/null/v8"

type BookingUnit struct {
	Id     string      `json:"id"`
	Ticket null.String `json:"ticket"`
}
