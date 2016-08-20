package memory

import (
	"hash/fnv"
	"sort"
	"time"
)

type multiStorage struct {
	storages []*storage
}

// NewMultiStorage creates new multi memory storage with n storages inside
func NewMultiStorage(n uint) *multiStorage {
	if n < 1 {
		n = 1
	}
	s := &multiStorage{}
	for i := uint(0); i < n; i++ {
		s.storages = append(s.storages, NewStorage())
	}
	return s
}

func (s *multiStorage) getStorage(key string) *storage {
	h := fnv.New32a()
	h.Write([]byte(key))
	n := int(h.Sum32()) % len(s.storages)
	return s.storages[n]
}

// Keys returns list of all keys
func (s *multiStorage) Keys() []string {
	var keys []string
	for _, storage := range s.storages {
		keys = append(keys, storage.Keys()...)
	}
	sort.Strings(keys)
	return keys
}

// TTL returns ttl of specified key. Error will occur if key doesn't exist.
func (s *multiStorage) TTL(key string) (time.Duration, error) {
	return s.getStorage(key).TTL(key)
}

// Get value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *multiStorage) Get(key string) (string, error) {
	return s.getStorage(key).Get(key)
}

// Set value of specified key with ttl. Use zero duration if key should exist forever.
// Error will occur if key already exists.
func (s *multiStorage) Set(key, value string, ttl time.Duration) error {
	return s.getStorage(key).Set(key, value, ttl)
}

// Update value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *multiStorage) Update(key, value string) error {
	return s.getStorage(key).Update(key, value)
}

// Delete specified key. Error will occur if key doesn't exist. It works for any key type.
func (s *multiStorage) Delete(key string) error {
	return s.getStorage(key).Delete(key)
}

// HashCreate creates new hash with specified key and ttl. Use zero duration if key should exist forever.
func (s *multiStorage) HashCreate(key string, ttl time.Duration) error {
	return s.getStorage(key).HashCreate(key, ttl)
}

// HashGet returns value of specified field of key.
// Error will occur if key or field doesn't exist or key type is not hash.
func (s *multiStorage) HashGet(key, field string) (string, error) {
	return s.getStorage(key).HashGet(key, field)
}

// HashGetAll returns all hash values of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *multiStorage) HashGetAll(key string) (map[string]string, error) {
	return s.getStorage(key).HashGetAll(key)
}

// HashSet sets field value of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *multiStorage) HashSet(key, field, value string) error {
	return s.getStorage(key).HashSet(key, field, value)
}

// HashDelete deletes field from hash. Error will occur if key doesn't exist or key type is not hash.
func (s *multiStorage) HashDelete(key, field string) error {
	return s.getStorage(key).HashDelete(key, field)
}

// HashLen returns count of hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *multiStorage) HashLen(key string) (int, error) {
	return s.getStorage(key).HashLen(key)
}

// HashKeys returns list of all hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *multiStorage) HashKeys(key string) ([]string, error) {
	return s.getStorage(key).HashKeys(key)
}

// ListCreate creates new list with specified key and ttl. Use zero duration if key should exist forever.
func (s *multiStorage) ListCreate(key string, ttl time.Duration) error {
	return s.getStorage(key).ListCreate(key, ttl)
}

// ListLeftPop pops value from the list beginning.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *multiStorage) ListLeftPop(key string) (string, error) {
	return s.getStorage(key).ListLeftPop(key)
}

// ListRightPop pops value from the list ending.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *multiStorage) ListRightPop(key string) (string, error) {
	return s.getStorage(key).ListRightPop(key)
}

// ListLeftPush adds value to the list beginning. Error will occur if key doesn't exist or key type is not list.
func (s *multiStorage) ListLeftPush(key, value string) error {
	return s.getStorage(key).ListLeftPush(key, value)
}

// ListRightPush adds value to the list ending. Error will occur if key doesn't exist or key type is not list.
func (s *multiStorage) ListRightPush(key, value string) error {
	return s.getStorage(key).ListRightPush(key, value)
}

// ListLen returns count of elements in the list. Error will occur if key doesn't exist or key type is not list.
func (s *multiStorage) ListLen(key string) (int, error) {
	return s.getStorage(key).ListLen(key)
}

// ListRange returns list of elements from the list from start to stop index.
// Error will occur if key doesn't exist or key type is not list.
func (s *multiStorage) ListRange(key string, start, stop int) ([]string, error) {
	return s.getStorage(key).ListRange(key, start, stop)
}
