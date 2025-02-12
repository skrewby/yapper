package models

import (
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skrewby/yapper/database"
	"github.com/skrewby/yapper/types"
)

type Users struct {
	db *pgxpool.Pool
}

func NewUsersModel(db *pgxpool.Pool) *Users {
	return &Users{
		db,
	}
}

func (u Users) CreateUser(email string, display_name string, hash string) error {
	query := `
		INSERT INTO users(email, display_name, password)
		VALUES (@email, @displayName, @hashPassword)
	`
	args := pgx.NamedArgs{
		"email":        email,
		"displayName":  display_name,
		"hashPassword": hash,
	}
	err := database.Run(u.db, query, args)
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"uc_email\" (SQLSTATE 23505)" {
			return types.CreateUserError{
				Msg:   "User with this email already exists",
				Field: "email",
			}
		}
	}

	return err
}

func (u Users) UpdateUser(user *types.User) error {
	query := `
		UPDATE users
		SET display_name = @display_name, active = @active
		WHERE id = @id
	`
	args := pgx.NamedArgs{
		"id":           user.Id,
		"display_name": user.Name,
		"active":       user.Active,
	}
	err := database.Run(u.db, query, args)

	return err
}

func (u Users) GetUser(id int) (*types.User, error) {
	query := `
		SELECT id, email, display_name, active, create_date::text, last_updated::text
		FROM users
		WHERE id = @id
	`
	args := pgx.NamedArgs{
		"id": id,
	}
	var user types.User
	err := database.SelectOne(u.db, &user, query, args)
	if err != nil {
		slog.Error("Get user", slog.String("Error", err.Error()))
	}

	return &user, err
}

func (u Users) GetUserByEmail(email string) (*types.User, error) {
	query := `
		SELECT id, email, display_name, active, create_date::text, last_updated::text
		FROM users
		WHERE email = @email
	`
	args := pgx.NamedArgs{
		"email": email,
	}
	var user types.User
	err := database.SelectOne(u.db, &user, query, args)
	if err != nil {
		slog.Error("Get user", slog.String("Error", err.Error()))
	}

	return &user, err
}

func (u Users) GetUserHashedPassword(id int) (*string, error) {
	query := `
		SELECT password
		FROM users
		WHERE id = @id
	`
	args := pgx.NamedArgs{
		"id": id,
	}
	var hash string
	err := database.SelectOne(u.db, &hash, query, args)
	if err != nil {
		slog.Error("Get user hashed password", slog.String("Error", err.Error()))
	}

	return &hash, err
}

func (u Users) GetAllUsers() ([]*types.User, error) {
	query := `
		SELECT id, email, display_name, active, create_date::text, last_updated::text
		FROM users
		ORDER BY display_name
	`
	var users []*types.User
	err := database.Select(u.db, &users, query, pgx.NamedArgs{})
	if err != nil {
		slog.Error("Get all users", slog.String("Error", err.Error()))
	}

	return users, err
}
