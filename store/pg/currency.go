package pg

import (
	"context"
	"github.com/google/uuid"
	"ventrata_task/dto"
	models "ventrata_task/store/generated/sqlBoiler"
)

type CurrencyPgRepo struct {
	db *DB
}

func NewCurrencyPgRepo(db *DB) *CurrencyPgRepo {
	return &CurrencyPgRepo{db}
}

func (r *CurrencyPgRepo) GetByCode(ctx context.Context, code dto.CurrencyCode) (*dto.Currency, error) {

	m, err := models.Currencies(models.CurrencyWhere.Code.EQ(code.String())).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	return &dto.Currency{
		Id:   id,
		Code: dto.CurrencyCode(m.Code),
	}, nil
}
