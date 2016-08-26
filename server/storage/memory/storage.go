package memory

import (
	"container/list"
	"sort"
	"sync"
	"time"

	commonStorage "github.com/Barberrrry/jcache/server/storage"
	"github.com/hashicorp/golang-lru/simplelru"
)

type storage struct {
	lru *simplelru.LRU
	mu  sync.RWMutex
}

// NewStorage creates new memory storage
func NewStorage(size int, gcInterval time.Duration) (*storage, error) {
	cache, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}

	s := &storage{lru: cache}

	go s.gc(gcInterval)

	return s, nil
}

func (s *storage) gc(interval time.Duration) {
	for _ = range time.Tick(interval) {
		var deleteKeys []interface{}
		s.mu.RLock()
		for _, key := range s.lru.Keys() {
			if raw, exists := s.lru.Get(key); exists {
				if item, castOk := raw.(*commonStorage.Item); castOk && !item.IsAlive() {
					deleteKeys = append(deleteKeys, key)
				}
			}
		}
		s.mu.RUnlock()
		if len(deleteKeys) > 0 {
			for _, key := range deleteKeys {
				s.mu.Lock()
				s.lru.Remove(key)
				s.mu.Unlock()
			}
		}
	}
}

func (s *storage) getItem(key string) (*commonStorage.Item, error) {
	if raw, exists := s.lru.Get(key); exists {
		if item, castOk := raw.(*commonStorage.Item); castOk && item.IsAlive() {
			return item, nil
		}
	}
	return nil, commonStorage.KeyNotExistsError
}

func (s *storage) addItem(key string, item *commonStorage.Item) {
	s.lru.Add(key, item)
}

func (s *storage) removeItem(key string) {
	s.lru.Remove(key)
}

func (s *storage) getHash(key string, createIfNotExist bool) (commonStorage.Hash, error) {
	item, err := s.getItem(key)
	if err != nil {
		if !createIfNotExist {
			return nil, err
		}
		item = commonStorage.NewItem(make(commonStorage.Hash), 0)
		s.addItem(key, item)
	}
	hash, err := item.CastHash()
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (s *storage) getList(key string, createIfNotExist bool) (*list.List, error) {
	item, err := s.getItem(key)
	if err != nil {
		if !createIfNotExist {
			return nil, err
		}
		item = commonStorage.NewItem(list.New(), 0)
		s.addItem(key, item)
	}
	list, err := item.CastList()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Keys returns list of all keys
func (s *storage) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, s.lru.Len())
	for _, key := range s.lru.Keys() {
		if item, err := s.getItem(key.(string)); err == nil && item.IsAlive() {
			keys = append(keys, key.(string))
		}
	}
	sort.Strings(keys)
	return keys
}

// Expire sets new key ttl
func (s *storage) Expire(key string, ttl uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, err := s.getItem(key)
	if err != nil {
		return err
	}

	item.SetTTL(ttl)
	return nil
}

// Get value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *storage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, err := s.getItem(key)
	if err != nil {
		return "", err
	}
	return item.CastString()
}

// Set value of specified key with ttl. Use zero ttl if key should exist forever.
// Error will occur if key already exists.
func (s *storage) Set(key, value string, ttl uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, _ := s.getItem(key)
	if item != nil {
		return commonStorage.KeyAlreadyExistsError
	}

	s.addItem(key, commonStorage.NewItem(value, ttl))
	return nil
}

// Update value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *storage) Update(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, err := s.getItem(key)
	if err != nil {
		return err
	}

	item.Value = value
	return nil
}

// Delete specified key. Error will occur if key doesn't exist. It works for any key type.
func (s *storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.getItem(key)
	if err != nil {
		return err
	}
	s.removeItem(key)
	return nil
}

// HashCreate creates new hash with specified key and ttl. Use zero ttl if key should exist forever.
func (s *storage) HashCreate(key string, ttl uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, _ := s.getItem(key)
	if item != nil {
		return commonStorage.KeyAlreadyExistsError
	}
	s.addItem(key, commonStorage.NewItem(make(commonStorage.Hash), ttl))
	return nil
}

// HashGet returns value of specified field of key.
// Error will occur if key or field doesn't exist or key type is not hash.
func (s *storage) HashGet(key, field string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash, err := s.getHash(key, false)
	if err != nil {
		return "", err
	}
	return hash.GetValue(field)
}

// HashGetAll returns all hash values of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashGetAll(key string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getHash(key, false)
}

// HashSet sets field value of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashSet(key, field, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, err := s.getHash(key, true)
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

	hash, err := s.getHash(key, false)
	if err != nil {
		return err
	}
	_, err = hash.GetValue(field)
	if err != nil {
		return err
	}
	delete(hash, field)
	return nil
}

// HashLen returns count of hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashLen(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash, err := s.getHash(key, false)
	if err != nil {
		return 0, err
	}
	return len(hash), nil
}

// HashKeys returns list of all hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashKeys(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash, err := s.getHash(key, false)
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

// ListCreate creates new list with specified key and ttl. Use zero ttl if key should exist forever.
func (s *storage) ListCreate(key string, ttl uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, _ := s.getItem(key)
	if item != nil {
		return commonStorage.KeyAlreadyExistsError
	}
	s.addItem(key, commonStorage.NewItem(list.New(), ttl))
	return nil
}

// ListLeftPop pops value from the list beginning.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListLeftPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key, false)
	if err != nil {
		return "", err
	}

	if e := list.Front(); e != nil {
		list.Remove(e)
		return e.Value.(string), nil
	}
	return "", commonStorage.ListEmptyError
}

// ListRightPop pops value from the list ending.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListRightPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key, false)
	if err != nil {
		return "", err
	}

	if e := list.Back(); e != nil {
		list.Remove(e)
		return e.Value.(string), nil
	}
	return "", commonStorage.ListEmptyError
}

// ListLeftPush adds value to the list beginning. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLeftPush(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, err := s.getList(key, true)
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

	list, err := s.getList(key, true)
	if err != nil {
		return err
	}

	list.PushBack(value)
	return nil
}

// ListLen returns count of elements in the list. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLen(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list, err := s.getList(key, false)
	if err != nil {
		return 0, err
	}

	return list.Len(), nil
}

// ListRange returns list of elements from the list from start to stop index.
// Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list, err := s.getList(key, false)
	if err != nil {
		return nil, err
	}

	var values []string
	e := list.Front()
	i := 0
	for {
		if e == nil || i > stop {
			break
		}
		if i >= start {
			values = append(values, e.Value.(string))
		}
		e = e.Next()
		i++
	}

	return values, nil
}
