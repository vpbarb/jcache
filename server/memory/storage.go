package memory

import (
	"container/list"
	"fmt"
	"sort"
	"sync"
	"time"
)

type storage struct {
	elements map[string]*element
	mu       sync.Mutex
}

// NewStorage creates new memory storage
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

func (s *storage) getList(key string) (*list.List, error) {
	element, err := s.getElement(key)
	if err != nil {
		return nil, err
	}
	list, err := element.castList()
	if err != nil {
		return nil, err
	}
	return list, nil
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

// Keys returns list of all keys
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

// TTL returns ttl of specified key. Error will occur if key doesn't exist.
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

// Get value of specified key. Error will occur if key doesn't exist or key type is not string.
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

// Set value of specified key with ttl. Use zero duration if key should exist forever.
// Error will occur if key already exists.
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

// Update value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *storage) Update(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	element, err := s.getElement(key)
	if err != nil {
		return err
	}

	element.value = value
	return nil
}

// Delete specified key. Error will occur if key doesn't exist. It works for any key type.
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

// HashCreate creates new hash with specified key and ttl. Use zero duration if key should exist forever.
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

// HashGet returns value of specified field of key.
// Error will occur if key or field doesn't exist or key type is not hash.
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

// HashGetAll returns all hash values of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashGetAll(key string) (map[string]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// HashSet sets field value of specified key. Error will occur if key doesn't exist or key type is not hash.
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

// HashDelete deletes field from hash. Error will occur if key doesn't exist or key type is not hash.
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

// HashLen returns count of hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key)
	if err != nil {
		return 0, err
	}

	return len(hash), nil
}

// HashKeys returns list of all hash fields. Error will occur if key doesn't exist or key type is not hash.
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

// ListCreate creates new list with specified key and ttl. Use zero duration if key should exist forever.
func (s *storage) ListCreate(key string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	element, _ := s.getElement(key)
	if element != nil {
		return fmt.Errorf(`Key "%s" already exists`, key)
	}
	s.elements[key] = s.createElement(key, list.New(), ttl)
	return nil
}

// ListLeftPop pops value from the list beginning.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListLeftPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key)
	if err != nil {
		return "", err
	}

	if e := list.Front(); e != nil {
		list.Remove(e)
		return e.Value.(string), nil
	}
	return "", fmt.Errorf(`List "%s" is empty`, key)
}

// ListRightPop pops value from the list ending.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListRightPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key)
	if err != nil {
		return "", err
	}

	if e := list.Back(); e != nil {
		list.Remove(e)
		return e.Value.(string), nil
	}
	return "", fmt.Errorf(`List "%s" is empty`, key)
}

// ListLeftPush adds value to the list beginning. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLeftPush(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key)
	if err != nil {
		return err
	}

	list.PushFront(value)
	return nil
}

// ListRightPush adds value to the list ending. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListRightPush(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key)
	if err != nil {
		return err
	}

	list.PushBack(value)
	return nil
}

// ListLen returns count of elements in the list. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key)
	if err != nil {
		return 0, err
	}

	return list.Len(), nil
}
