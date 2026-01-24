package database

import (
	"context"
	"errors"
	"os"
	"time"

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

	dbConnStr := os.Getenv("DATABASE_URL")

	if dbConnStr == "" {
		msg := "invalid database connection string.\n please specify env variable DATABASE_URL."
		logger.L().Fatalln(msg)
	}

	//TODO: implement pool.Close() on app exist or panic
	pool, err := pgxpool.New(context.Background(), dbConnStr)
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
