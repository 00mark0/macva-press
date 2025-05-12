package utils

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(uuidStr, fieldName string) (pgtype.UUID, error) {
	uuidBytes, err := uuid.Parse(uuidStr)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid %s format: %w", fieldName, err)
	}
	return pgtype.UUID{
		Bytes: uuidBytes,
		Valid: true,
	}, nil
}
