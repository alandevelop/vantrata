package store

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"ventrata_task/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunPgMigrations runs Postgres migrations
func RunPgMigrations() error {
	cfg := config.Get()

	if cfg.PgUrl == "" {
		return errors.New("No cfg.PgURL provided")
	}

	p := []tryPathParam{
		{"file://store/pg/migrations", cfg.PgUrl},
		{"file://../store/pg/migrations", cfg.PgUrl},
		{"file://../../store/pg/migrations", cfg.PgUrl},
	}

	m, err := tryPath(p)

	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

type tryPathParam struct {
	source      string
	databaseURL string
}

func tryPath(p []tryPathParam) (*migrate.Migrate, error) {

	var m *migrate.Migrate
	var err error

	for _, p := range p {
		m = nil
		err = nil

		m, err = migrate.New(
			p.source,
			p.databaseURL,
		)

		if err != nil {
			continue
		}

		return m, err
	}

	return m, err
}
