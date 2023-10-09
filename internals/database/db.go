package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	MaxIdleConnections int
	MaxOpenConnections int
	MaxIdleTime        time.Duration
}

func OpenDB(uri string, config DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetConnMaxIdleTime(config.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
