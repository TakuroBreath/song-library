package migrator

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
)

func Migrate(user, password, host, port, dbname string, log *slog.Logger) error {
	m, err := migrate.New("file://migrations", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname))
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no changes to apply")

			return nil
		}

		return err
	}

	log.Info("migrations applied")

	return nil
}
