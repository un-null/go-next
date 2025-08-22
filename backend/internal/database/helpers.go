// database/helpers.go
package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToPgtype(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

func PgtypeToUUID(pgUUID pgtype.UUID) uuid.UUID {
	if !pgUUID.Valid {
		return uuid.Nil
	}
	return pgUUID.Bytes
}

func TimeToPgtype(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

func PgtypeToTime(pgTime pgtype.Timestamptz) time.Time {
	if !pgTime.Valid {
		return time.Time{}
	}
	return pgTime.Time
}

func Int32ToPgtype(i int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: i,
		Valid: true,
	}
}

func PgtypeToInt32(pgInt pgtype.Int4) int32 {
	return pgInt.Int32
}

func NumericToFloat64(num pgtype.Numeric) float64 {
	if !num.Valid {
		return 0
	}

	f, _ := num.Float64Value()
	return f.Float64
}
