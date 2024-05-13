package store

import (
	"context"
	"github.com/google/uuid"
	"ventrata_task/dto"
)

type ProductRepo interface {
	Create(context.Context, *dto.Product) error
	Get(context.Context, uuid.UUID) (*dto.Product, error)
	Exists(context.Context, uuid.UUID) (bool, error)
	List(context.Context) ([]*dto.Product, error)
}

type AvailabilityRepo interface {
	Get(context.Context, uuid.UUID) (*dto.Availability, error)
	GetWhere(context.Context, *dto.AvailabilityFilter) ([]dto.Availability, error)
	Update(context.Context, *dto.Availability) error
	Seed(context.Context, []dto.Availability) error
}

type BookingsRepo interface {
	Get(context.Context, uuid.UUID) (*dto.Booking, error)
	Create(context.Context, *dto.Booking) error
	Update(context.Context, *dto.Booking) error
	List(context.Context) ([]*dto.Booking, error)
	CreateTickets(context.Context, uuid.UUID) ([]dto.BookingUnit, error)
}

type CurrencyRepo interface {
	GetByCode(context.Context, dto.CurrencyCode) (*dto.Currency, error)
}
