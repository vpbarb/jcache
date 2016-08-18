package memory

import (
	"fmt"
	"sync"
	//"container/list"
	"sort"
	"time"
)

type storage struct {
	elements map[string]*element
	mu       sync.Mutex
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

func (s *storage) getHash(key string) (hash, error) {
	element, err := s.getElement(key)
	if err != nil {
		return nil, err
	}
	hash, err := element.castHash()
	if err != nil {
		return nil, err
	}
	return hash, nil
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
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := make([]string, 0, len(s.elements))
	for key := range s.elements {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *storage) TTL(key string) (time.Duration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

	element, err := s.getElement(key)
	if err != nil {
		return "", err
	}
	value, err := element.castString()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *storage) Set(key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	element, _ := s.getElement(key)
	if element != nil {
		return fmt.Errorf(`Key "%s" already exists`, key)
	}

	s.elements[key] = s.createElement(key, value, ttl)
	return nil
}

// todo
func (s *storage) Update(key, value string) error {
	return nil
}

func (s *storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.getElement(key)
	if err != nil {
		return err
	}
	delete(s.elements, key)
	return nil
}

func (s *storage) HashCreate(key string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	element, _ := s.getElement(key)
	if element != nil {
		return fmt.Errorf(`Key "%s" already exists`, key)
	}
	s.elements[key] = s.createElement(key, make(hash), ttl)
	return nil
}

func (s *storage) HashGet(key, field string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return "", err
	}
	value, err := hash.getValue(field)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *storage) HashGetAll(key string) (map[string]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (s *storage) HashSet(key, field, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return err
	}
	hash[field] = value
	return nil
}

func (s *storage) HashDelete(key, field string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return err
	}
	_, err = hash.getValue(field)
	if err != nil {
		return err
	}
	delete(hash, field)
	return nil
}

func (s *storage) HashLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return 0, err
	}

	return len(hash), nil
}

func (s *storage) HashKeys(key string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(hash))
	for key := range hash {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

// todo
func (s *storage) ListLeftPop(key string) (string, error) {
	return "", nil
}

// todo
func (s *storage) ListRightPop(key string) (string, error) {
	return "", nil
}

// todo
func (s *storage) ListLeftPush(key, value string, ttl time.Duration) error {
	return nil
}

// todo
func (s *storage) ListRightPush(key, value string, ttl time.Duration) error {
	return nil
}

// todo
func (s *storage) ListSet(key string, index int, value string, ttl time.Duration) error {
	return nil
}

// todo
func (s *storage) ListIndex(key string, index int) (string, error) {
	return "", nil
}

// todo
func (s *storage) ListLen(key string) (int, error) {
	return 0, nil
}

// todo
func (s *storage) ListDelete(key string, count int, value string) error {
	return nil
}

// todo
func (s *storage) ListRange(key string, start, stop int) ([]string, error) {
	return []string{}, nil
}
