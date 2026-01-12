package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"dubai-auto/internal/config"
)

func Init(cfg *config.Config) *pgxpool.Pool {
	connectionString := buildConnectionString(cfg)
	dbConfig, err := pgxpool.ParseConfig(connectionString)

	if err != nil {
		log.Fatalf("Unable to parse database config💊: %v\n", err)
	}

	dbConfig.MaxConns = 200
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		log.Fatalf("failed to create connection poolpool🏊: %v\n", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		panic(fmt.Sprintf("Could not ping postgres🫙 database: %v", err))
	}

	log.Println("Database 🥳 connection pool initialized successfully ✅")
	return pool
}

func buildConnectionString(cfg *config.Config) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DB_USER, cfg.DB_PASSWORD,
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME,
	)
}
