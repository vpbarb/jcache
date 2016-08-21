package storage

import (
	"container/list"
	"fmt"
	"time"
)

type Item struct {
	Value      interface{}
	ExpireTime time.Time
}

func NewItem(value interface{}, ttl uint64) *Item {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &Item{
		Value:      value,
		ExpireTime: expireTime,
	}
}

func (i *Item) CastString() (string, error) {
	if value, ok := i.Value.(string); ok {
		return value, nil
	} else {
		return "", fmt.Errorf(`Key type is not string`)
	}
}

func (i *Item) CastHash() (Hash, error) {
	if hash, ok := i.Value.(Hash); ok {
		return hash, nil
	} else {
		return nil, fmt.Errorf(`Key type is not hash`)
	}
}

func (i *Item) CastList() (*list.List, error) {
	if list, ok := i.Value.(*list.List); ok {
		return list, nil
	} else {
		return nil, fmt.Errorf(`Key type is not list`)
	}
}

func (i *Item) IsAlive() bool {
	return i.ExpireTime.IsZero() || i.ExpireTime.After(time.Now())
}

func (i *Item) GetTTL() uint64 {
	if i.ExpireTime.IsZero() {
		return uint64(0)
	}
	d := i.ExpireTime.Sub(time.Now())
	return uint64(d.Seconds())
}
