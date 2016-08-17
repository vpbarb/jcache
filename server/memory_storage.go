package server

import (
	"fmt"
	"sync"
)

type Hash map[string]string

type List []string

type memoryStorage struct {
	values map[string]interface{}
	m      sync.Mutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		values: make(map[string]interface{}),
	}
}

func (s *memoryStorage) Keys() ([]string, error) {
	return []string{}, nil
}

func (s *memoryStorage) findValue(key string) (interface{}, error) {
	if value, found := s.values[key]; found {
		return value, nil
	}
	return nil, fmt.Errorf(`Key "%s" is not exists`, key)
}

func (s *memoryStorage) Get(key string) (string, error) {
	value, err := s.findValue(key)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (s *memoryStorage) Set(key, value string) error {
	s.values[key] = value
	return nil
}

func (s *memoryStorage) Update(key, value string) error {
	return nil
}

func (s *memoryStorage) Delete(key string) error {
	_, err := s.findValue(key)
	if err != nil {
		return err
	}
	delete(s.values, key)
	return nil
}

func (s *memoryStorage) HashGet(key, field string) (string, error) {
	return "", nil
}

func (s *memoryStorage) HashSet(key, field, value string) error {
	return nil
}

func (s *memoryStorage) HashUpdate(key, field, value string) error {
	return nil
}

func (s *memoryStorage) HashDelete(key, field string) error {
	return nil
}

func (s *memoryStorage) HashLen(key string) (int, error) {
	return 0, nil
}

func (s *memoryStorage) HashKeys(key string) ([]string, error) {
	return []string{}, nil
}

func (s *memoryStorage) ListLeftPop(key string) (string, error) {
	return "", nil
}

func (s *memoryStorage) ListRightPop(key string) (string, error) {
	return "", nil
}

func (s *memoryStorage) ListLeftPush(key, value string) error {
	return nil
}

func (s *memoryStorage) ListRightPush(key, value string) error {
	return nil
}

func (s *memoryStorage) ListSet(key string, index int, value string) error {
	return nil
}

func (s *memoryStorage) ListIndex(key string, index int) (string, error) {
	return "", nil
}

func (s *memoryStorage) ListLen(key string) (int, error) {
	return 0, nil
}

func (s *memoryStorage) ListDelete(key string, count int, value string) error {
	return nil
}

func (s *memoryStorage) ListRange(key string, start, stop int) ([]string, error) {
	return []string{}, nil
}
