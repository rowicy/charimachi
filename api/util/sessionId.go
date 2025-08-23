package util

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateSessionID() string {
	value := uuid.New().String()
	fmt.Printf("Generated session_id: %s\n", value) // TODO デバッグ用
	return value
}

var SessionIDResponse = map[string]ORSGeometry{}
