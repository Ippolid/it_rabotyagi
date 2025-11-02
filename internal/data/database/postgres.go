package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"it_rabotyagi/internal/logger"
)

// DB представляет собой обертку над пулом соединений с базой данных.
type DB struct {
	Pool *pgxpool.Pool
}

// NewPostgresConnection создает и возвращает новый экземпляр DB.
func NewPostgresConnection(url string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close закрывает все соединения в пуле.
func (db *DB) Close() {
	//log.Println("Closing database connection pool.")
	logger.Info("Closing database connection pool.")
	db.Pool.Close()
}
