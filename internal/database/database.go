package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/arsenh/DevLore/internal/config"
	"github.com/arsenh/DevLore/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DbContext struct {
	pool *pgxpool.Pool
}

func NewDbContext() *DbContext {
	return &DbContext{}
}

func (d *DbContext) MustConnect() {
	err := godotenv.Load() // TODO: consider to add this in configs
	if err != nil {
		logger.L().Fatalln("failed to load .env file.")
		panic(err)
	}

	dbConnectionURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	if dbConnectionURL == "" {
		msg := "invalid database connection string.\n please specify env variable DATABASE_URL."
		logger.L().Fatalln(msg)
	}

	pool, err := pgxpool.New(context.Background(), dbConnectionURL)
	if err != nil {
		logger.L().Infoln("faild connect to database.")
		panic(err)
	}
	d.pool = pool

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		panic(err)
	}

	logger.L().Infoln("connected to the database successfully.")
	config.DatabaseConnectionString = dbConnectionURL
}

func (d *DbContext) Close() {
	d.pool.Close()
}

func (d *DbContext) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, sql, args...)
}

func (d *DbContext) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return d.pool.QueryRow(ctx, sql, args...)
}

func (d *DbContext) QueryRows(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if d.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}
	return d.pool.Query(ctx, sql, args...)
}
