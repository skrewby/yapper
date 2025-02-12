package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetDateTimeStr(ts pgtype.Timestamptz) string {
	return ts.Time.Format(time.RFC822)
}
