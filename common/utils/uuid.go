package utils

import "github.com/google/uuid"

var uuidGenerator = uuid.New()

func NewUUID() string {
	return uuidGenerator.String()
}
