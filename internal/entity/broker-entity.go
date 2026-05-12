package entity

import "github.com/google/uuid"

type Broker struct {
	ID   uuid.UUID
	Name string
}

func NewBroker(name string) (Broker, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Broker{}, err
	}
	return Broker{
		ID:   id,
		Name: name,
	}, nil
}
