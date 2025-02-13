package models

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skrewby/yapper/database"
	"github.com/skrewby/yapper/types"
)

type Threads struct {
	db *pgxpool.Pool
}

func NewThreadsModel(db *pgxpool.Pool) *Threads {
	return &Threads{
		db,
	}
}

func (m Threads) CreateThread(title string, userId int) error {
	query := `
		INSERT INTO threads(title, author)
		VALUES (@title, @author)
	`
	args := pgx.NamedArgs{
		"title":  title,
		"author": userId,
	}
	err := database.Run(m.db, query, args)

	return err
}

func (m Threads) GetAllThreads() ([]*types.Thread, error) {
	query := `
		SELECT t.id, t.title, t.create_date::text, u.id AS "author.id", u.email AS "author.email", u.display_name AS "author.display_name"
		FROM threads t
		JOIN users u ON u.id = t.author
		ORDER BY t.create_date  DESC 
	`
	var threads []*types.Thread
	err := database.Select(m.db, &threads, query, pgx.NamedArgs{})
	if err != nil {
		slog.Error("Get all threads", slog.String("Error", err.Error()))
	}

	return threads, err
}
