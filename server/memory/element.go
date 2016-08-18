package memory

import (
	"fmt"
	"time"
)

type element struct {
	value      interface{}
	expireTime time.Time
}

func (e *element) castHash() (hash, error) {
	if hash, ok := e.value.(hash); ok {
		return hash, nil
	} else {
		return nil, fmt.Errorf(`Key type is not hash`)
	}
}

func (e *element) castString() (string, error) {
	if value, ok := e.value.(string); ok {
		return value, nil
	} else {
		return "", fmt.Errorf(`Key type is not string`)
	}
}
