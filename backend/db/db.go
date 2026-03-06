package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB(connString string) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	}

	if pool != nil {
		if err := pool.Ping(ctx); err != nil {
			log.Printf("Database not reachable: %v\n", err)
		} else {
			log.Println("Connected to PostgreSQL")
		}
	}

	Pool = pool
}
