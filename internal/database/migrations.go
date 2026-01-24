package database

import (
	"github.com/arsenh/DevLore/internal/config"
	"github.com/arsenh/DevLore/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	m, err := migrate.New(
		"file://internal/migrations",
		config.DatabaseConnectionString,
	)

	if err != nil {
		logger.L().Fatalln("failed to create migrations object", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.L().Fatalln("failed to up migrations", err)
	}

	logger.L().Println("database migrations applied successfully")
}
