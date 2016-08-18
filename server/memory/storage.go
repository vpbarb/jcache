package memory

import (
	"fmt"
	"sync"
	//"container/list"
	"sort"
	"time"
)

type element struct {
	value      interface{}
	expireTime time.Time
}

type storage struct {
	elements map[string]*element
	m        sync.Mutex
}

func NewStorage() *storage {
	return &storage{
		elements: make(map[string]*element),
	}
}

func (s *storage) getElement(key string) (*element, error) {
	if element, found := s.elements[key]; found {
		if element.expireTime.IsZero() || element.expireTime.After(time.Now()) {
			return element, nil
		}
		delete(s.elements, key)
	}
	return nil, fmt.Errorf(`Key "%s" does not exist`, key)
}

func (s *storage) createElement(key string, value interface{}, ttl time.Duration) *element {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(ttl)
	}
	return &element{
		value:      value,
		expireTime: expireTime,
	}
}

func (s *storage) Keys() []string {
	keys := make([]string, 0, len(s.elements))
	for key := range s.elements {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *storage) TTL(key string) (time.Duration, error) {
	element, err := s.getElement(key)
	if err != nil {
		return time.Duration(0), err
	}
	if element.expireTime.IsZero() {
		return time.Duration(0), nil
	}
	return element.expireTime.Sub(time.Now()), nil
}

func (s *storage) Get(key string) (string, error) {
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

func (s *storage) Set(key, value string, ttl time.Duration) error {
	s.elements[key] = s.createElement(key, value, ttl)
	return nil
}

func (s *storage) Update(key, value string) error {
	return nil
}

func (s *storage) Delete(key string) error {
	_, err := s.getElement(key)
	if err != nil {
		return err
	}
	delete(s.elements, key)
	return nil
}

func (s *storage) HashCreate(key string, ttl time.Duration) error {
	element, _ := s.getElement(key)
	if element != nil {
		return fmt.Errorf(`Hash with key "%s" already exists`, key)
	}
	s.elements[key] = s.createElement(key, make(map[string]string), ttl)
	return nil
}

func (s *storage) HashGet(key, field string) (string, error) {
	element, err := s.getElement(key)
	if err != nil {
		return "", err
	}
	if hash, ok := element.value.(map[string]string); ok {
		if value, found := hash[field]; found {
			return value, nil
		} else {
			return "", fmt.Errorf(`Field "%s" does not exist`, field)
		}
	} else {
		return "", fmt.Errorf(`Key "%s" is not hash`, key)
	}
}

func (s *storage) HashGetAll(key string) (map[string]string, error) {
	element, err := s.getElement(key)
	if err != nil {
		return nil, err
	}
	if hash, ok := element.value.(map[string]string); ok {
		return hash, nil
	} else {
		return nil, fmt.Errorf(`Key "%s" is not hash`, key)
	}
}

func (s *storage) HashSet(key, field, value string) error {
	element, err := s.getElement(key)
	if err != nil {
		return err
	}
	if hash, ok := element.value.(map[string]string); ok {
		hash[field] = value
		return nil
	} else {
		return fmt.Errorf(`Key "%s" is not hash`, key)
	}
}

func (s *storage) HashDelete(key, field string) error {
	return nil
}

func (s *storage) HashLen(key string) (int, error) {
	return 0, nil
}

func (s *storage) HashKeys(key string) ([]string, error) {
	return []string{}, nil
}

func (s *storage) ListLeftPop(key string) (string, error) {
	return "", nil
}

func (s *storage) ListRightPop(key string) (string, error) {
	return "", nil
}

func (s *storage) ListLeftPush(key, value string, ttl time.Duration) error {
	return nil
}

func (s *storage) ListRightPush(key, value string, ttl time.Duration) error {
	return nil
}

func (s *storage) ListSet(key string, index int, value string, ttl time.Duration) error {
	return nil
}

func (s *storage) ListIndex(key string, index int) (string, error) {
	return "", nil
}

func (s *storage) ListLen(key string) (int, error) {
	return 0, nil
}

func (s *storage) ListDelete(key string, count int, value string) error {
	return nil
}

func (s *storage) ListRange(key string, start, stop int) ([]string, error) {
	return []string{}, nil
}
