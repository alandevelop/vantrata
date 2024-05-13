package dto

import (
	"github.com/google/uuid"
	"time"
)

type Availability struct {
	Id        uuid.UUID
	ProductId uuid.UUID
	LocalDate time.Time
	Vacancies int
	Status    AvailabilityStatus
	Available bool
}

type AvailabilityWeb struct {
	Id        string `json:"id,omitempty"`
	ProductId string `json:"-"`
	LocalDate string `json:"localDate,omitempty"`
	Vacancies int    `json:"vacancies,omitempty"`
	Status    string `json:"status,omitempty"`
	Available bool   `json:"available,omitempty"`
}

func (a *Availability) ToWeb() AvailabilityWeb {
	return AvailabilityWeb{
		Id:        a.Id.String(),
		ProductId: a.ProductId.String(),
		LocalDate: a.LocalDate.Format(time.DateOnly),
		Vacancies: a.Vacancies,
		Status:    a.Status.String(),
		Available: a.Available,
	}
}

type AvailabilityFilter struct {
	ProductId     uuid.UUID `json:"productId,omitempty" param:"id" validate:"required"`
	LocalDate     time.Time `json:"localDate" param:"localDate"`
	LocalDateFrom time.Time `json:"localDateFrom" param:"localDateFrom"`
	LocalDateTo   time.Time `json:"localDateTo" param:"localDateTo"`
}

type AvailabilityStatus string

const (
	AvailabilityStatusAvailable AvailabilityStatus = "AVAILABLE"
	AvailabilityStatusSoldOut                      = "SOLD_OUT"
)

func (s AvailabilityStatus) String() string {
	return string(s)
}
