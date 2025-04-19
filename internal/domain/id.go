package domain

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID = primitive.ObjectID

func NewID(id string) (ID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ID{}, fmt.Errorf("invalid ID: %w", err)
	}
	return oid, nil
}

func NewAutoID() ID {
	return primitive.NewObjectID()
}
