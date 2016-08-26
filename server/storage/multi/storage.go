package multi

import (
	"hash/fnv"
	"sort"

	commonStorage "github.com/Barberrrry/jcache/server/storage"
)

// Multi storage allow combine different storage types in one.
// Distribution by storages is normal and based on key hash sum.
type storage struct {
	storages []commonStorage.Storage
}

// NewStorage creates new multi memory storage with n storages inside
func NewStorage(storages ...commonStorage.Storage) *storage {
	return &storage{storages: storages}
}

func (s *storage) AddStorage(storage commonStorage.Storage) {
	s.storages = append(s.storages, storage)
}

func (s *storage) getStorage(key string) commonStorage.Storage {
	h := fnv.New32a()
	h.Write([]byte(key))
	n := int(h.Sum32()) % len(s.storages)
	return s.storages[n]
}

// Keys returns list of all keys
func (s *storage) Keys() []string {
	var keys []string
	for _, storage := range s.storages {
		keys = append(keys, storage.Keys()...)
	}
	sort.Strings(keys)
	return keys
}

// Expire sets new key ttl
func (s *storage) Expire(key string, ttl uint64) error {
	return s.getStorage(key).Expire(key, ttl)
}

// Get value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *storage) Get(key string) (string, error) {
	return s.getStorage(key).Get(key)
}

// Set value of specified key with ttl. Use zero duration if key should exist forever.
// Error will occur if key already exists.
func (s *storage) Set(key, value string, ttl uint64) error {
	return s.getStorage(key).Set(key, value, ttl)
}

// Update value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *storage) Update(key, value string) error {
	return s.getStorage(key).Update(key, value)
}

// Delete specified key. Error will occur if key doesn't exist. It works for any key type.
func (s *storage) Delete(key string) error {
	return s.getStorage(key).Delete(key)
}

// HashCreate creates new hash with specified key and ttl. Use zero duration if key should exist forever.
func (s *storage) HashCreate(key string, ttl uint64) error {
	return s.getStorage(key).HashCreate(key, ttl)
}

// HashGet returns value of specified field of key.
// Error will occur if key or field doesn't exist or key type is not hash.
func (s *storage) HashGet(key, field string) (string, error) {
	return s.getStorage(key).HashGet(key, field)
}

// HashGetAll returns all hash values of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashGetAll(key string) (map[string]string, error) {
	return s.getStorage(key).HashGetAll(key)
}

// HashSet sets field value of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashSet(key, field, value string) error {
	return s.getStorage(key).HashSet(key, field, value)
}

// HashDelete deletes field from hash. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashDelete(key, field string) error {
	return s.getStorage(key).HashDelete(key, field)
}

// HashLen returns count of hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashLen(key string) (int, error) {
	return s.getStorage(key).HashLen(key)
}

// HashKeys returns list of all hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashKeys(key string) ([]string, error) {
	return s.getStorage(key).HashKeys(key)
}

// ListCreate creates new list with specified key and ttl. Use zero duration if key should exist forever.
func (s *storage) ListCreate(key string, ttl uint64) error {
	return s.getStorage(key).ListCreate(key, ttl)
}

// ListLeftPop pops value from the list beginning.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListLeftPop(key string) (string, error) {
	return s.getStorage(key).ListLeftPop(key)
}

// ListRightPop pops value from the list ending.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListRightPop(key string) (string, error) {
	return s.getStorage(key).ListRightPop(key)
}

// ListLeftPush adds value to the list beginning. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLeftPush(key, value string) error {
	return s.getStorage(key).ListLeftPush(key, value)
}

// ListRightPush adds value to the list ending. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListRightPush(key, value string) error {
	return s.getStorage(key).ListRightPush(key, value)
}

// ListLen returns count of elements in the list. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLen(key string) (int, error) {
	return s.getStorage(key).ListLen(key)
}

// ListRange returns list of elements from the list from start to stop index.
// Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListRange(key string, start, stop int) ([]string, error) {
	return s.getStorage(key).ListRange(key, start, stop)
}
