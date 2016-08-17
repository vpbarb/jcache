package server

import (
	"fmt"
	"sync"
	//"container/list"
	"time"
)

type element struct {
	value      interface{}
	expireTime time.Time
}

type Hash map[string]string

type memoryStorage struct {
	elements map[string]*element
	m        sync.Mutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		elements: make(map[string]*element),
	}
}

func (s *memoryStorage) getElement(key string) (*element, error) {
	if element, found := s.elements[key]; found {
		if element.expireTime.IsZero() || element.expireTime.After(time.Now()) {
			return element, nil
		}
		delete(s.elements, key)
	}
	return nil, fmt.Errorf(`Key "%s" is not exists`, key)
}

func (s *memoryStorage) Keys() []string {
	keys := make([]string, 0, len(s.elements))
	for key := range s.elements {
		keys = append(keys, key)
	}
	return keys
}

func (s *memoryStorage) TTL(key string) (time.Duration, error) {
	element, err := s.getElement(key)
	if err != nil {
		return time.Duration(0), err
	}
	if element.expireTime.IsZero() {
		return time.Duration(0), nil
	}
	return element.expireTime.Sub(time.Now()), nil
}

func (s *memoryStorage) Get(key string) (string, error) {
	element, err := s.getElement(key)
	if err != nil {
		return "", err
	}
	return element.value.(string), nil
}

func (s *memoryStorage) Set(key, value string, ttl time.Duration) error {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(ttl)
	}
	s.elements[key] = &element{
		value:      value,
		expireTime: expireTime,
	}
	return nil
}

func (s *memoryStorage) Update(key, value string) error {
	return nil
}

func (s *memoryStorage) Delete(key string) error {
	_, err := s.getElement(key)
	if err != nil {
		return err
	}
	delete(s.elements, key)
	return nil
}

func (s *memoryStorage) HashGet(key, field string) (string, error) {
	return "", nil
}

func (s *memoryStorage) HashSet(key, field, value string, ttl time.Duration) error {
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

func (s *memoryStorage) ListLeftPush(key, value string, ttl time.Duration) error {
	return nil
}

func (s *memoryStorage) ListRightPush(key, value string, ttl time.Duration) error {
	return nil
}

func (s *memoryStorage) ListSet(key string, index int, value string, ttl time.Duration) error {
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
