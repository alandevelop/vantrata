package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	models "ventrata_task/store/generated/sqlBoiler"
	"ventrata_task/store/pg"
)

type Store struct {
	DB *pg.DB

	Product      ProductRepo
	Availability AvailabilityRepo
	Bookings     BookingsRepo
	Currency     CurrencyRepo
}

// New creates new store
func New(ctx context.Context) (*Store, error) {
	// connect to Postgres
	pgDB, err := pg.Dial()
	if err != nil {
		return nil, errors.Wrap(err, "pgdb.Dial failed")
	}

	// Run Postgres migrations
	if pgDB != nil {
		log.Println("Running PostgreSQL migrations...")
		if err := RunPgMigrations(); err != nil {
			return nil, errors.Wrap(err, "runPgMigrations failed")
		}
	}

	var store Store

	// Init Postgres repositories
	if pgDB != nil {
		store.DB = pgDB
		store.Availability = pg.NewAvailabilityRepo(pgDB)
		store.Bookings = pg.NewBookingPgRepo(pgDB)
		store.Currency = pg.NewCurrencyPgRepo(pgDB)
		store.Product = pg.NewProductPgRepo(pgDB)

		Seed(pgDB.DB)
	}

	return &store, nil
}

func Seed(db *sql.DB) {
	fmt.Println("running seeds")

	ctx := context.Background()

	c := models.Currency{
		ID:   "0de138b8-a5e9-4944-baa9-2eb692de14ea",
		Code: "EUR",
	}

	err := c.Upsert(ctx, db, false, nil, boil.Infer(), boil.Infer())
	if err != nil {
		fmt.Println(err)
	}

	p := models.Product{
		ID:         "364d965c-ad1d-4450-a554-cbbaf536ffa3",
		Name:       "test_product",
		Capacity:   10,
		Price:      1000,
		CurrencyID: "0de138b8-a5e9-4944-baa9-2eb692de14ea",
	}

	err = p.Upsert(ctx, db, false, nil, boil.Infer(), boil.Infer())
	if err != nil {
		fmt.Println(err)
	}
}
