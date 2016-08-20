package storage

import (
	"time"
)

type Storage interface {
	Keys() []string
	TTL(key string) (time.Duration, error)
	Get(key string) (string, error)
	Set(key, value string, ttl time.Duration) error
	Update(key, value string) error
	Delete(key string) error
	HashCreate(key string, ttl time.Duration) error
	HashGet(key, field string) (string, error)
	HashGetAll(key string) (map[string]string, error)
	HashSet(key, field, value string) error
	HashDelete(key, field string) error
	HashLen(key string) (int, error)
	HashKeys(key string) ([]string, error)
	ListCreate(key string, ttl time.Duration) error
	ListLeftPop(key string) (string, error)
	ListRightPop(key string) (string, error)
	ListLeftPush(key, value string) error
	ListRightPush(key, value string) error
	ListLen(key string) (int, error)
	ListRange(key string, start, stop int) ([]string, error)
}
