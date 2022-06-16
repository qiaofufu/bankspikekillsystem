package utils

import "github.com/google/uuid"

func GetUUID() (string, error) {
	UUID, err := uuid.NewRandom()
	return UUID.String(), err
}
