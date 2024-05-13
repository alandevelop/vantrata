package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
	"ventrata_task/dto"
	"ventrata_task/lib/types"
	"ventrata_task/store"
)

type AvailabilitySvc struct {
	store *store.Store
}

var seedLock = make(map[string]int)

func NewAvailabilitySvc(store *store.Store) *AvailabilitySvc {
	return &AvailabilitySvc{store}
}

func (s AvailabilitySvc) List(ctx context.Context, filter *dto.AvailabilityFilter) ([]dto.Availability, error) {
	avs, err := s.store.Availability.GetWhere(ctx, filter)

	if err != nil {
		return nil, s.handleError("vc.Availability.List", "Availability not found", err)
	}

	if avs != nil && len(avs) == 0 {
		return nil, errors.Wrap(types.ErrNotFound, "Availability not found")
	}

	return avs, nil
}

func (s AvailabilitySvc) SeedForYear(ctx context.Context, productId uuid.UUID) error {
	product, err := s.store.Product.Get(ctx, productId)

	if err != nil {
		return s.handleError("vc.Availability.SeedForYear", "Availability not found", err)
	}

	// check lock
	if _, ok := seedLock[productId.String()]; ok {
		return nil
	}

	// set lock
	seedLock[productId.String()]++

	today := time.Now().Truncate(24 * time.Hour)
	end := today.AddDate(0, 0, 366)

	f := &dto.AvailabilityFilter{
		ProductId:     productId,
		LocalDateFrom: today,
		LocalDateTo:   end,
	}

	avs, err := s.store.Availability.GetWhere(ctx, f)
	if err != nil {
		return errors.Wrap(types.ErrNotFound, "Availability not found")
	}

	var result []dto.Availability

	if len(avs) == 0 {
		result = s.newDaysSeed(product)
	} else {
		result = s.idempotentSeed(product, avs)
	}

	if len(result) == 0 {
		return nil
	}

	err = s.store.Availability.Seed(ctx, result)
	if err != nil {
		errors.Wrap(err, "svc.Availability.SeedForYear")
	}

	// release lock
	delete(seedLock, productId.String())

	return nil
}

func (s AvailabilitySvc) idempotentSeed(product *dto.Product, avs []dto.Availability) []dto.Availability {
	d := 0 // overall year
	newDate := time.Now()
	result := make([]dto.Availability, 0, 366)

out:
	for i, a := range avs {

		var nextDate time.Time
		if len(avs) > i+1 {
			nextDate = avs[i+1].LocalDate
		}

		for d < 366 {

			// not to overwrite current
			if newDate.Year() == a.LocalDate.Year() && newDate.YearDay() == a.LocalDate.YearDay() {
				d++
				newDate = newDate.AddDate(0, 0, 1)
				continue out
			}

			// not to overwrite  next date
			if !nextDate.IsZero() && (newDate.Year() == nextDate.Year() && newDate.YearDay() == nextDate.YearDay()) {
				d++
				newDate = newDate.AddDate(0, 0, 1)
				continue out
			}

			avNew := dto.Availability{
				ProductId: product.Id,
				LocalDate: newDate,
				Vacancies: product.Capacity,
				Status:    dto.AvailabilityStatusAvailable,
				Available: true,
			}

			result = append(result, avNew)

			d++
			newDate = newDate.AddDate(0, 0, 1)
		}

	}

	return result
}

func (s AvailabilitySvc) newDaysSeed(product *dto.Product) []dto.Availability {
	newDate := time.Now()
	result := make([]dto.Availability, 0, 366)

	for i := 0; i < 366; i++ {
		avNew := dto.Availability{
			ProductId: product.Id,
			LocalDate: newDate,
			Vacancies: product.Capacity,
			Status:    dto.AvailabilityStatusAvailable,
			Available: true,
		}

		result = append(result, avNew)

		newDate = newDate.AddDate(0, 0, 1)
	}

	return result
}

func (s AvailabilitySvc) handleError(msg string, msgNotFound string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(types.ErrNotFound, msgNotFound)
	} else {
		return errors.Wrap(err, msg)
	}
}
