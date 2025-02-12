package types

import "github.com/jackc/pgx/v5/pgtype"

type Auth struct {
	Token string `json:"token"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	Email    string `json:"email"`
	Name     string `json:"display_name" db:"display_name"`
	Password string `json:"password"`
}

type User struct {
	Id       int                `json:"id"`
	Email    string             `json:"email"`
	Name     string             `json:"display_name" db:"display_name"`
	Password string             `json:"password,omitempty"`
	Active   *bool              `json:"active"`
	Created  pgtype.Timestamptz `json:"created" db:"create_date"`
	Updated  pgtype.Timestamptz `json:"updated" db:"last_updated"`
}

type JWTUser struct {
	Email string `json:"email"`
	Name  string `json:"display_name" db:"display_name"`
}

type JWTContext struct {
	User JWTUser `json:"user"`
}

type JWT struct {
	Context JWTContext `json:"context"`
}
