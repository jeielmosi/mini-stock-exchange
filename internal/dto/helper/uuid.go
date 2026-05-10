package dto_helper

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

func EncodeUUID(id uuid.UUID) (string, error) {
	b, err := id.MarshalBinary()
	if err != nil {
		return "", err
	}
	str := base64.RawURLEncoding.EncodeToString(b)
	return str, nil
}

func DecodeUUID(str string) (uuid.UUID, error) {
	b, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = id.UnmarshalBinary(b)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func DecodeUUIDv7(str string) (uuid.UUID, error) {
	id, err := DecodeUUID(str)
	if err != nil {
		return uuid.Nil, err
	}
	err = IsValidUUIDv7(id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func IsValidUUIDv7(id uuid.UUID) error {
	if id.Version() != 7 {
		return fmt.Errorf("invalid UUID version, expected 7 but got %d", id.Version())
	}
	return nil
}
