package pg

import (
	"context"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"ventrata_task/dto"
	models "ventrata_task/store/generated/sqlBoiler"
)

type ProductPgRepo struct {
	db *DB
}

func NewProductPgRepo(db *DB) *ProductPgRepo {
	return &ProductPgRepo{db: db}
}

func (r *ProductPgRepo) Create(ctx context.Context, p *dto.Product) error {
	m := &models.Product{
		ID:         p.Id.String(),
		Name:       p.Name,
		Capacity:   p.Capacity,
		Price:      p.Price,
		CurrencyID: p.Currency.Id.String(),
	}

	return m.Insert(ctx, r.db, boil.Infer())
}

func (r *ProductPgRepo) Get(ctx context.Context, pid uuid.UUID) (*dto.Product, error) {
	m, err := models.Products(
		models.ProductWhere.ID.EQ(pid.String()),
		qm.Load(models.ProductRels.Currency),
	).One(ctx, r.db)

	if err != nil {
		return nil, err
	}

	cid, err := uuid.Parse(m.R.Currency.ID)
	if err != nil {
		return nil, err
	}

	return &dto.Product{
		Id:       pid,
		Name:     m.Name,
		Capacity: m.Capacity,
		Price:    m.Price,
		Currency: &dto.Currency{
			Id:   cid,
			Code: dto.CurrencyCode(m.R.Currency.Code),
		},
	}, nil
}

func (r *ProductPgRepo) Exists(ctx context.Context, pid uuid.UUID) (bool, error) {
	return models.Products(
		models.ProductWhere.ID.EQ(pid.String()),
	).Exists(ctx, r.db)
}

func (r *ProductPgRepo) List(ctx context.Context) ([]*dto.Product, error) {
	mods, err := models.Products(
		qm.Load(models.ProductRels.Currency),
	).All(ctx, r.db)

	if err != nil {
		return nil, err
	}

	res := make([]*dto.Product, 0, len(mods))

	for _, m := range mods {

		pid, err := uuid.Parse(m.ID)
		if err != nil {
			return nil, err
		}

		cid, err := uuid.Parse(m.R.Currency.ID)
		if err != nil {
			return nil, err
		}

		p := &dto.Product{
			Id:       pid,
			Name:     m.Name,
			Capacity: m.Capacity,
			Price:    m.Price,
			Currency: &dto.Currency{
				Id:   cid,
				Code: dto.CurrencyCode(m.R.Currency.Code),
			},
		}

		res = append(res, p)
	}

	return res, nil
}
