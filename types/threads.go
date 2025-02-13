package types

import "github.com/jackc/pgx/v5/pgtype"

type Thread struct {
	Id      int                `json:"id"`
	Title   string             `json:"email"`
	Author  User               `json:"author"`
	Created pgtype.Timestamptz `json:"created" db:"create_date"`
}
