package dto

import (
	"github.com/google/uuid"
)

type Booking struct {
	Id             uuid.UUID
	ProductId      uuid.UUID
	AvailabilityId uuid.UUID
	Status         BookingStatus
	Price          int
	Currency       *Currency
	Units          []BookingUnit
}

type CreateBookingRequest struct {
	ProductId      string `json:"productId" validate:"required" form:"productId"`
	AvailabilityId string `json:"availabilityId" validate:"required" form:"availabilityId"`
	Units          int    `json:"units" validate:"required" form:"units"`
}

type BookingWeb struct {
	Id             string        `json:"id"`
	ProductId      string        `json:"productId"`
	AvailabilityId string        `json:"availabilityId"`
	Status         string        `json:"status"`
	Price          int           `json:"price,omitempty"`
	Currency       string        `json:"currency,omitempty"`
	Units          []BookingUnit `json:"units,omitempty"`
}

func (b Booking) ToWeb() BookingWeb {
	return BookingWeb{
		Id:             b.Id.String(),
		ProductId:      b.ProductId.String(),
		AvailabilityId: b.AvailabilityId.String(),
		Status:         b.Status.String(),
		Price:          b.Price,
		Currency:       b.Currency.Code.String(),
		Units:          b.Units,
	}
}

type BookingStatus string

const (
	BookingStatusReserved  BookingStatus = "RESERVED"
	BookingStatusConfirmed               = "CONFIRMED"
)

func (s BookingStatus) String() string {
	return string(s)
}
