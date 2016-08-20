package memory

import (
	"container/list"
	"fmt"
	"time"
)

type item struct {
	value      interface{}
	expireTime time.Time
}

func (i *item) castString() (string, error) {
	if value, ok := i.value.(string); ok {
		return value, nil
	} else {
		return "", fmt.Errorf(`Key type is not string`)
	}
}

func (i *item) castHash() (hash, error) {
	if hash, ok := i.value.(hash); ok {
		return hash, nil
	} else {
		return nil, fmt.Errorf(`Key type is not hash`)
	}
}

func (i *item) castList() (*list.List, error) {
	if list, ok := i.value.(*list.List); ok {
		return list, nil
	} else {
		return nil, fmt.Errorf(`Key type is not list`)
	}
}
