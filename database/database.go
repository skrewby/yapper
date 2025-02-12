package database

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skrewby/yapper/utils"
)

func SelectOne(db *pgxpool.Pool, data interface{}, query string, args pgx.NamedArgs) error {
	rows, err := db.Query(context.Background(), query, args)
	if err != nil {
		return err
	}

	if err = pgxscan.ScanOne(data, rows); err != nil {
		return err
	}

	return nil
}

func Select(db *pgxpool.Pool, data interface{}, query string, args pgx.NamedArgs) error {
	err := pgxscan.Select(context.Background(), db, data, query, args)

	return err
}

func Run(db *pgxpool.Pool, query string, args pgx.NamedArgs) error {
	_, err := db.Exec(context.Background(), query, args)

	return err
}

func ConnectDatabase(env utils.Environment) (*pgxpool.Pool, error) {
	// Because it's currently running inside a container, then the location of db is @db:5432
	db_user := env.DatabaseCredentials.User
	db_pass := env.DatabaseCredentials.Pass
	db_host := "db"
	db_name := env.DatabaseCredentials.Name
	db_port := "5432"

	db_url := "postgres://" + db_user + ":" + db_pass + "@" + db_host + ":" + db_port + "/" + db_name
	conn, err := pgxpool.New(context.Background(), db_url)
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool (%s)", err)
	}

	return conn, nil
}
