package pg

import (
	"context"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"ventrata_task/dto"
	models "ventrata_task/store/generated/sqlBoiler"
)

type BookingPgRepo struct {
	db *DB
}

func NewBookingPgRepo(db *DB) *BookingPgRepo {
	return &BookingPgRepo{db: db}
}

func (r *BookingPgRepo) Get(ctx context.Context, pid uuid.UUID) (*dto.Booking, error) {
	m, err := models.Bookings(
		models.BookingWhere.ID.EQ(pid.String()),
		qm.Load(models.BookingRels.Currency),
		qm.Load(models.BookingRels.BookingUnits),
	).One(ctx, r.db)

	if err != nil {
		return nil, err
	}

	bid, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	aid, err := uuid.Parse(m.AvailabilityID)
	if err != nil {
		return nil, err
	}

	cid, err := uuid.Parse(m.R.Currency.ID)
	if err != nil {
		return nil, err
	}

	res := dto.Booking{
		Id:             bid,
		ProductId:      pid,
		AvailabilityId: aid,
		Status:         dto.BookingStatus(m.Status.String()),
		Price:          m.Price,
		Currency: &dto.Currency{
			Id:   cid,
			Code: dto.CurrencyCode(m.R.Currency.Code),
		},
		Units: []dto.BookingUnit{},
	}

	if m.R != nil && m.R.BookingUnits != nil {
		for _, u := range m.R.BookingUnits {
			temp := dto.BookingUnit{
				Id:     u.ID,
				Ticket: u.Ticket,
			}
			res.Units = append(res.Units, temp)
		}
	}

	return &res, nil
}

func (r *BookingPgRepo) Create(ctx context.Context, b *dto.Booking) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	m := models.Booking{
		ID:             b.Id.String(),
		ProductID:      b.ProductId.String(),
		AvailabilityID: b.AvailabilityId.String(),
		Status:         models.BookingStatus(b.Status.String()),
		Price:          b.Price,
		CurrencyID:     b.Currency.Id.String(),
	}

	err = m.Insert(ctx, tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := 0; i < len(b.Units); i++ {
		u := models.BookingUnit{
			ID:        uuid.NewString(),
			BookingID: b.Id.String(),
			Ticket:    null.String{},
		}

		err = u.Insert(ctx, tx, boil.Infer())
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *BookingPgRepo) Update(ctx context.Context, b *dto.Booking) error {

	m, err := models.Bookings(
		models.BookingWhere.ID.EQ(b.Id.String()),
		qm.Load(models.BookingRels.Currency),
	).One(ctx, r.db)

	if err != nil {
		return err
	}

	m.Status = models.BookingStatus(b.Status.String())

	_, err = m.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *BookingPgRepo) List(ctx context.Context) ([]*dto.Booking, error) {
	mods, err := models.Bookings(
		qm.Load(models.BookingRels.Currency),
	).All(ctx, r.db)

	if err != nil {
		return nil, err
	}

	res := make([]*dto.Booking, 0, len(mods))

	for _, m := range mods {
		bid, err := uuid.Parse(m.ID)
		if err != nil {
			return nil, err
		}

		aid, err := uuid.Parse(m.AvailabilityID)
		if err != nil {
			return nil, err
		}

		pid, err := uuid.Parse(m.ProductID)
		if err != nil {
			return nil, err
		}

		cid, err := uuid.Parse(m.R.Currency.ID)
		if err != nil {
			return nil, err
		}

		b := &dto.Booking{
			Id:             bid,
			ProductId:      pid,
			AvailabilityId: aid,
			Status:         dto.BookingStatus(m.Status.String()),
			Price:          m.Price,
			Currency: &dto.Currency{
				Id:   cid,
				Code: dto.CurrencyCode(m.R.Currency.Code),
			},
		}

		res = append(res, b)
	}

	return res, nil
}

func (r *BookingPgRepo) CreateTickets(ctx context.Context, bookingId uuid.UUID) ([]dto.BookingUnit, error) {

	units, err := models.BookingUnits(models.BookingUnitWhere.BookingID.EQ(bookingId.String())).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var res []dto.BookingUnit

	for _, u := range units {
		u.Ticket = null.String{
			String: uuid.NewString(),
			Valid:  true,
		}

		_, err = u.Update(ctx, tx, boil.Infer())
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		res = append(res, dto.BookingUnit{
			Id:     u.ID,
			Ticket: u.Ticket,
		})
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return res, nil
}
