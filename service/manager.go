package service

import (
	"github.com/pkg/errors"
	"ventrata_task/store"
)

type Manager struct {
	Product      ProductService
	Availability AvailabilityService
	Bookings     BookingsService
}

// NewManager creates new service manager
func NewManager(store *store.Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("No store provided")
	}
	return &Manager{
		Product:      NewProductSvc(store),
		Availability: NewAvailabilitySvc(store),
		Bookings:     NewBookingsSvc(store),
	}, nil
}
