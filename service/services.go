package service

import (
	"context"
	"github.com/google/uuid"
	"ventrata_task/dto"
)

type ProductService interface {
	Create(context.Context, *dto.Product) (*dto.Product, error)
	Get(context.Context, uuid.UUID) (*dto.Product, error)
	List(context.Context) ([]*dto.Product, error)
}

type AvailabilityService interface {
	List(context.Context, *dto.AvailabilityFilter) ([]dto.Availability, error)
	SeedForYear(context.Context, uuid.UUID) error
}

type BookingsService interface {
	GetBooking(context.Context, uuid.UUID) (*dto.Booking, error)
	CreateBooking(context.Context, *dto.CreateBookingRequest) (*dto.Booking, error)
	Confirm(context.Context, uuid.UUID) (*dto.Booking, error)
	List(context.Context) ([]*dto.Booking, error)
}
