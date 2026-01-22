package database

import (
	"context"
	"os"

	"github.com/arsenh/DevLore/internal/logger"
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
		logger.L().Fatalln("invalid database connection string.")
	}

	//TODO: implement pool.Close() on app exist or panic
	pool, err := pgxpool.New(context.Background(), dbConnStr)
	if err != nil {
		logger.L().Infoln("faild connect to database.")
		panic(err)
	}
	d.pool = pool
	logger.L().Infoln("connected to the database successfully.")
}
