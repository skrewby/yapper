package models_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		slog.Error("Could not construct pool", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	err = pool.Client.Ping()
	if err != nil {
		slog.Error("Could not connect to Docker", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
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
		slog.Error("Could not start resource", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	slog.Info("Connecting to database on URL: ", slog.String("URL", databaseUrl))

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = pgxpool.New(context.Background(), databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	}); err != nil {
		slog.Error("Could not connect to docker", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			slog.Error("Could not purge resource", slog.String("Error", err.Error()))
			os.Exit(1)
		}
	}()

	slog.Info("Running migrations")
	sqlDB := stdlib.OpenDBFromPool(db)
	gooseProvider, err := goose.NewProvider(
		goose.DialectPostgres,
		sqlDB,
		os.DirFS("../sql"),
	)
	if err != nil {
		slog.Error("Could not create a goose provider to the db", slog.String("Error", err.Error()))
		os.Exit(1)
	}
	_, err = gooseProvider.Up(context.Background())
	if err != nil {
		slog.Error("Error applying migrations", slog.String("Error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Running tests")

	m.Run()
}
