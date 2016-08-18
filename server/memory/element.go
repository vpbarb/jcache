package memory

import (
	"container/list"
	"fmt"
	"time"
)

type element struct {
	value      interface{}
	expireTime time.Time
}

func (e *element) castString() (string, error) {
	if value, ok := e.value.(string); ok {
		return value, nil
	} else {
		return "", fmt.Errorf(`Key type is not string`)
	}
}

func (e *element) castHash() (hash, error) {
	if hash, ok := e.value.(hash); ok {
		return hash, nil
	} else {
		return nil, fmt.Errorf(`Key type is not hash`)
	}
}

func (e *element) castList() (*list.List, error) {
	if list, ok := e.value.(*list.List); ok {
		return list, nil
	} else {
		return nil, fmt.Errorf(`Key type is not list`)
	}
}
