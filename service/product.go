package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"ventrata_task/dto"
	"ventrata_task/lib/types"
	"ventrata_task/store"
)

type ProductSvc struct {
	store *store.Store
}

func NewProductSvc(store *store.Store) *ProductSvc {
	return &ProductSvc{store}
}

func (s ProductSvc) Create(ctx context.Context, product *dto.Product) (*dto.Product, error) {
	product.Id = uuid.New()
	product.Currency = &dto.CurrencyEUR

	err := s.store.Product.Create(ctx, product)

	if err != nil {
		return nil, errors.Wrap(err, "svc.Product.Create")
	}

	return product, nil
}

func (s ProductSvc) Get(ctx context.Context, productId uuid.UUID) (*dto.Product, error) {
	product, err := s.store.Product.Get(ctx, productId)

	if err != nil {
		return nil, s.handleError("svc.Product.Get", err)
	}

	return product, nil
}

func (s ProductSvc) List(ctx context.Context) ([]*dto.Product, error) {
	products, err := s.store.Product.List(ctx)

	if err != nil {
		return nil, s.handleError("svc.Product.List", err)
	}

	return products, nil
}

func (s ProductSvc) handleError(msg string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(types.ErrNotFound, "Product not found")
	} else {
		return errors.Wrap(err, msg)
	}
}
