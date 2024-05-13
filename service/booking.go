package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"ventrata_task/dto"
	"ventrata_task/lib/types"
	"ventrata_task/store"
)

type BookingsSvc struct {
	store *store.Store
}

func NewBookingsSvc(store *store.Store) *BookingsSvc {
	return &BookingsSvc{store: store}
}

func (s *BookingsSvc) GetBooking(ctx context.Context, bookingId uuid.UUID) (*dto.Booking, error) {
	bk, err := s.store.Bookings.Get(ctx, bookingId)

	if err != nil {
		return nil, s.handleError("svc.Booking.GetBooking", "Booking not found", err)
	}

	return bk, nil
}

func (s *BookingsSvc) CreateBooking(ctx context.Context, request *dto.CreateBookingRequest) (*dto.Booking, error) {
	if request.Units == 0 {
		return nil, errors.Wrap(types.ErrNotFound, "Not enough units")
	}

	productId, err := uuid.Parse(request.ProductId)
	if err != nil {
		return nil, errors.Wrap(err, "svc.Booking.CreateBooking")
	}

	product, err := s.store.Product.Get(ctx, productId)
	if err != nil {
		return nil, s.handleError("svc.Booking.CreateBooking", "Product not found", err)
	}

	availabilityId, err := uuid.Parse(request.AvailabilityId)
	if err != nil {
		return nil, errors.Wrap(err, "svc.Booking.CreateBooking")
	}

	av, err := s.store.Availability.Get(ctx, availabilityId)
	if err != nil {
		return nil, s.handleError("svc.Booking.CreateBooking", "Availability not found", err)
	}

	if !av.Available || av.Vacancies < request.Units {
		return nil, errors.Wrap(types.ErrBadRequest, "not enough availability")
	}

	av.Vacancies = av.Vacancies - request.Units

	if av.Vacancies == 0 {
		av.Available = false
		av.Status = dto.AvailabilityStatusSoldOut
	}

	err = s.store.Availability.Update(ctx, av)
	if err != nil {
		return nil, errors.Wrap(err, "svc.Booking.CreateBooking")
	}

	curr, err := s.store.Currency.GetByCode(ctx, dto.CurrencyCodeEUR)
	if err != nil {
		return nil, errors.Wrap(err, "svc.Booking.CreateBooking")
	}

	bid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "svc.Booking.CreateBooking")
	}

	var units []dto.BookingUnit

	for i := 0; i < request.Units; i++ {
		t := dto.BookingUnit{
			Id:     uuid.NewString(),
			Ticket: null.String{},
		}

		units = append(units, t)
	}

	createBooking := &dto.Booking{
		Id:             bid,
		ProductId:      productId,
		AvailabilityId: availabilityId,
		Status:         dto.BookingStatusReserved,
		Price:          product.Price * request.Units,
		Currency:       curr,
		Units:          units,
	}

	err = s.store.Bookings.Create(ctx, createBooking)
	if err != nil {
		return nil, errors.Wrap(err, "svc.Booking.CreateBooking")
	}

	return createBooking, nil
}

func (s *BookingsSvc) Confirm(ctx context.Context, bookingId uuid.UUID) (*dto.Booking, error) {
	bk, err := s.store.Bookings.Get(ctx, bookingId)

	if err != nil {
		return nil, s.handleError("svc.Booking.Confirm", "Booking not found", err)
	}

	if bk.Status == dto.BookingStatusConfirmed {
		return nil, errors.Wrap(types.ErrDuplicateEntry, "booking is already confirmed")
	}

	bk.Status = dto.BookingStatusConfirmed

	err = s.store.Bookings.Update(ctx, bk)
	if err != nil {
		return nil, errors.Wrap(err, "svc.Booking.Confirm")
	}

	units, err := s.store.Bookings.CreateTickets(ctx, bk.Id)
	bk.Units = units

	return bk, nil
}

func (s *BookingsSvc) List(ctx context.Context) ([]*dto.Booking, error) {
	bookings, err := s.store.Bookings.List(ctx)

	if err != nil {
		return nil, s.handleError("svc.Booking.List", "Booking not found", err)
	}

	return bookings, nil
}

func (s BookingsSvc) handleError(msg string, msgNotFound string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(types.ErrNotFound, msgNotFound)
	} else {
		return errors.Wrap(err, msg)
	}
}
