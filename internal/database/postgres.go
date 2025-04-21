package database

import (
	"database/sql"
	"fmt"

	"github.com/alexey-shedrin/avito-test-task/internal/config"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func New(cfg *config.Config) (*sql.DB, error) {
	conStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	if err = goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	return db, nil
}

func IsUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
