package db

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DbPool *pgxpool.Pool

func InitStore() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	DbPool, err = pgxpool.New(context.Background(), dbPath)

	if err != nil {
		log.Fatal(err, "Unable to connect to database")
	}

	log.Println("Connected to database")
}

func AcqureConnection(ctx context.Context) (*pgxpool.Conn, error) {
	if DbPool == nil {
		return nil, errors.New("database pool not initialized")
	}

	conn, err := DbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
