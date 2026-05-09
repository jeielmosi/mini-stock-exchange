package dto_helper

import (
	"fmt"

	"github.com/google/uuid"
)

func IsValidUUIDv7(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	if id.Version() == 7 {
		return nil
	}
	return fmt.Errorf("invalid UUID version, expected 7 but got %d", id.Version())
}
