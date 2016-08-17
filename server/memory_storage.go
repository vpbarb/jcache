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
	return nil, fmt.Errorf(`Key "%s" does not exist`, key)
}

func (s *memoryStorage) createElement(key string, value interface{}, ttl time.Duration) *element {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(ttl)
	}
	return &element{
		value:      value,
		expireTime: expireTime,
	}
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
	if value, ok := element.value.(string); ok {
		return value, nil
	} else {
		return "", fmt.Errorf(`Key "%s" is not string`, key)
	}
}

func (s *memoryStorage) Set(key, value string, ttl time.Duration) error {
	s.elements[key] = s.createElement(key, value, ttl)
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

func (s *memoryStorage) HashCreate(key string, ttl time.Duration) error {
	element, _ := s.getElement(key)
	if element != nil {
		return fmt.Errorf(`Hash with key "%s" already exists`, key)
	}
	s.elements[key] = s.createElement(key, Hash{}, ttl)
	return nil
}

func (s *memoryStorage) HashGet(key, field string) (string, error) {
	element, err := s.getElement(key)
	if err != nil {
		return "", err
	}
	if hash, ok := element.value.(Hash); ok {
		if value, found := hash[field]; found {
			return value, nil
		} else {
			return "", fmt.Errorf(`Field "%s" does not exist`, field)
		}
	} else {
		return "", fmt.Errorf(`Key "%s" is not hash`, key)
	}
}

func (s *memoryStorage) HashGetAll(key string) (Hash, error) {
	element, err := s.getElement(key)
	if err != nil {
		return Hash{}, err
	}
	if hash, ok := element.value.(Hash); ok {
		return hash, nil
	} else {
		return Hash{}, fmt.Errorf(`Key "%s" is not hash`, key)
	}
}

func (s *memoryStorage) HashSet(key, field, value string) error {
	element, err := s.getElement(key)
	if err != nil {
		return err
	}
	if hash, ok := element.value.(Hash); ok {
		hash[field] = value
		return nil
	} else {
		return fmt.Errorf(`Key "%s" is not hash`, key)
	}
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
