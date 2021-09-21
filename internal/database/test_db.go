package database

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
)

const (
	migrationFilesPath = "../../migrations"
)

type CleanUp func() error

// createTestDB spins up a temporary docker database to run tests against
func createTestDB() (*pgxpool.Pool, CleanUp, error) {
	ctx := context.Background()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not connect to docker")
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not start resource")
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	connStr := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	// tell docker to hard kill the container in 120 seconds
	resource.Expire(120)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second

	var pgpool *pgxpool.Pool
	if err = pool.Retry(func() error {
		pgpool, err = pgxpool.Connect(ctx, connStr)
		if err != nil {
			return err
		}

		return pgpool.Ping(ctx)
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	if err := applyTestMigrations(connStr); err != nil {
		return nil, nil, err
	}

	cleanUp := func() error {
		return pool.Purge(resource)
	}

	return pgpool, cleanUp, nil
}

func applyTestMigrations(connStr string) error {
	path, err := filepath.Abs(migrationFilesPath)
	if err != nil {
		return errors.Wrap(err, "failed to create path to migrations")
	}

	m, err := migrate.New("file://"+path, connStr)
	if err != nil {
		return errors.Wrap(err, "failed to connect to test DB")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to run migrations")
	}

	return nil
}
