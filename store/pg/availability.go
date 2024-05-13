package pg

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"ventrata_task/dto"
	models "ventrata_task/store/generated/sqlBoiler"
)

type AvailabilityPgRepo struct {
	db *DB
}

func NewAvailabilityRepo(db *DB) *AvailabilityPgRepo {
	return &AvailabilityPgRepo{db: db}
}

func (r AvailabilityPgRepo) Seed(ctx context.Context, data []dto.Availability) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, d := range data {

		m := models.Availability{
			ID:        uuid.New().String(),
			ProductID: d.ProductId.String(),
			Localdate: d.LocalDate,
			Status:    models.AvailabilityStatus(d.Status.String()),
			Vacancies: d.Vacancies,
			Available: d.Available,
		}

		if err := m.Insert(ctx, tx, boil.Infer()); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r AvailabilityPgRepo) Get(ctx context.Context, id uuid.UUID) (*dto.Availability, error) {
	m, err := models.Availabilities(
		models.AvailabilityWhere.ID.EQ(id.String()),
		qm.Load(models.AvailabilityRels.Product),
	).One(ctx, r.db)

	if err != nil {
		return nil, err
	}

	aid, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	pid, err := uuid.Parse(m.R.Product.ID)
	if err != nil {
		return nil, err
	}

	return &dto.Availability{
		Id:        aid,
		ProductId: pid,
		LocalDate: m.Localdate,
		Vacancies: m.Vacancies,
		Status:    dto.AvailabilityStatus(m.Status.String()),
		Available: m.Available,
	}, nil
}

func (r AvailabilityPgRepo) GetWhere(ctx context.Context, f *dto.AvailabilityFilter) ([]dto.Availability, error) {

	var query []qm.QueryMod
	query = append(query, models.AvailabilityWhere.ProductID.EQ(f.ProductId.String()))

	if !f.LocalDate.IsZero() {
		query = append(query, models.AvailabilityWhere.Localdate.EQ(f.LocalDate))
	}

	if !f.LocalDateFrom.IsZero() && !f.LocalDateTo.IsZero() {
		query = append(query, qm.Where(fmt.Sprintf("%s BETWEEN ? AND ?", models.AvailabilityColumns.Localdate), f.LocalDateFrom, f.LocalDateTo))
	}

	query = append(query,
		qm.Load(models.AvailabilityRels.Product),
		qm.OrderBy(models.AvailabilityColumns.Localdate),
	)

	mods, err := models.Availabilities(
		query...,
	).All(ctx, r.db)

	if err != nil {
		return nil, err
	}

	res := make([]dto.Availability, 0, len(mods))

	for _, m := range mods {
		aid, err := uuid.Parse(m.ID)
		if err != nil {
			return nil, err
		}

		pid, err := uuid.Parse(m.R.Product.ID)
		if err != nil {
			return nil, err
		}

		a := dto.Availability{
			Id:        aid,
			ProductId: pid,
			LocalDate: m.Localdate,
			Vacancies: m.Vacancies,
			Status:    dto.AvailabilityStatus(m.Status.String()),
			Available: m.Available,
		}

		res = append(res, a)
	}

	return res, nil
}

func (r AvailabilityPgRepo) Update(ctx context.Context, a *dto.Availability) error {
	m := models.Availability{
		ID:        a.Id.String(),
		ProductID: a.ProductId.String(),
		Localdate: a.LocalDate,
		Status:    models.AvailabilityStatus(a.Status.String()),
		Vacancies: a.Vacancies,
		Available: a.Available,
	}

	_, err := m.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
