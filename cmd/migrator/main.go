package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	DatabaseDSN string
}

func loadConfig() (*Config, error) {
	viper.AutomaticEnv()

	c := &Config{
		DatabaseDSN: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			viper.GetString("DATABASE_USER"),
			viper.GetString("DATABASE_PASSWORD"),
			viper.GetString("DATABASE_HOST"),
			viper.GetString("DATABASE_PORT"),
			viper.GetString("DATABASE_DB"),
		),
	}

	return c, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "loading config"))
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(errors.Wrap(err, "creating logger"))
	}
	defer logger.Sync()

	path, err := filepath.Abs("./migrations/")
	if err != nil {
		logger.Error("unable to get migrations path", zap.Error(err))
	}

	m, err := migrate.New("file://"+path, cfg.DatabaseDSN)
	if err != nil {
		logger.Error("failed to connect to the database", zap.Error(err))
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error("failed to run migrations", zap.Error(err))
	}

	logger.Error("migrations successfully ran", zap.Error(err))
}
