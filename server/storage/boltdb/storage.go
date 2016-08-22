package boltdb

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sort"

	commonStorage "github.com/Barberrrry/jcache/server/storage"
	"github.com/boltdb/bolt"
)

var (
	bucketName        = []byte("bucket")
	notSupportedError = errors.New("Operation is not supported by BoltDB storage")
)

// Storage uses BoltDB as a persistent file-based storage.
// encoding/gob is used to encode/decode data structures to put them into BoltDB.
// Unfortunately container/list couldn't be used in a such way, so this storage doesn't support lists :(
// It may be implemented by custom list solution or by using some different encoder/decoder.
type storage struct {
	db *bolt.DB
}

func init() {
	gob.Register(commonStorage.Item{})
	gob.Register(commonStorage.Hash{})
}

func NewStorage(filePath string) (*storage, error) {
	db, err := bolt.Open(filePath, 0644, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot open Bolt file: %s", err)
	}

	return &storage{db: db}, nil
}

func (s *storage) getItem(tx *bolt.Tx, key string) (*commonStorage.Item, error) {
	bucket, err := tx.CreateBucketIfNotExists(bucketName)
	if err != nil {
		return nil, err
	}
	data := bucket.Get([]byte(key))
	if data == nil {
		return nil, fmt.Errorf(`Key "%s" does not exist`, key)
	}

	dec := gob.NewDecoder(bytes.NewBuffer(data))
	var item commonStorage.Item
	err = dec.Decode(&item)
	if err != nil {
		return nil, err
	}

	if item.IsAlive() {
		return &item, nil
	} else {
		bucket.Delete([]byte(key))
		return nil, fmt.Errorf(`Key "%s" does not exist`, key)
	}
}

func (s *storage) saveItem(tx *bolt.Tx, key string, item *commonStorage.Item) error {
	bucket, err := tx.CreateBucketIfNotExists(bucketName)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err = enc.Encode(item)
	if err != nil {
		return err
	}
	err = bucket.Put([]byte(key), buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) getHash(tx *bolt.Tx, key string) (commonStorage.Hash, error) {
	item, err := s.getItem(tx, key)
	if err != nil {
		return nil, err
	}
	hash, err := item.CastHash()
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// Keys returns list of all keys
func (s *storage) Keys() (keys []string) {
	s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			dec := gob.NewDecoder(bytes.NewBuffer(value))
			item := &commonStorage.Item{}
			err = dec.Decode(item)
			if err != nil {
				return err
			}

			if item.IsAlive() {
				keys = append(keys, string(key))
			}
		}
		return nil
	})
	sort.Strings(keys)
	return
}

// Get value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *storage) Get(key string) (value string, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		item, err := s.getItem(tx, key)
		if err != nil {
			return err
		}
		value, err = item.CastString()
		return err
	})
	return
}

// Set value of specified key with ttl. Use zero ttl if key should exist forever.
// Error will occur if key already exists.
func (s *storage) Set(key, value string, ttl uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		item, _ := s.getItem(tx, key)
		if item != nil {
			return fmt.Errorf(`Key "%s" already exists`, key)
		}

		item = commonStorage.NewItem(value, ttl)
		return s.saveItem(tx, key, item)
	})
}

// Update value of specified key. Error will occur if key doesn't exist or key type is not string.
func (s *storage) Update(key, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		item, err := s.getItem(tx, key)
		if err != nil {
			return err
		}

		item.Value = value
		return s.saveItem(tx, key, item)
	})
}

// Delete specified key. Error will occur if key doesn't exist. It works for any key type.
func (s *storage) Delete(key string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := s.getItem(tx, key)
		if err != nil {
			return err
		}

		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return bucket.Delete([]byte(key))
	})
}

// HashCreate creates new hash with specified key and ttl. Use zero ttl if key should exist forever.
func (s *storage) HashCreate(key string, ttl uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		item, _ := s.getItem(tx, key)
		if item != nil {
			return fmt.Errorf(`Key "%s" already exists`, key)
		}

		item = commonStorage.NewItem(make(commonStorage.Hash), ttl)
		return s.saveItem(tx, key, item)
	})
}

// HashGet returns value of specified field of key.
// Error will occur if key or field doesn't exist or key type is not hash.
func (s *storage) HashGet(key, field string) (value string, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		hash, err := s.getHash(tx, key)
		if err != nil {
			return err
		}
		value, err = hash.GetValue(field)
		return err
	})
	return
}

// HashGetAll returns all hash values of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashGetAll(key string) (hash map[string]string, err error) {
	err = s.db.Update(func(tx *bolt.Tx) (err error) {
		hash, err = s.getHash(tx, key)
		return err
	})
	return
}

// HashSet sets field value of specified key. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashSet(key, field, value string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		item, err := s.getItem(tx, key)
		if err != nil {
			return err
		}
		hash, err := item.CastHash()
		if err != nil {
			return err
		}
		hash[field] = value

		return s.saveItem(tx, key, item)
	})
}

// HashDelete deletes field from hash. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashDelete(key, field string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		item, err := s.getItem(tx, key)
		if err != nil {
			return err
		}
		hash, err := item.CastHash()
		if err != nil {
			return err
		}
		_, err = hash.GetValue(field)
		if err != nil {
			return err
		}
		delete(hash, field)
		return s.saveItem(tx, key, item)
	})
}

// HashLen returns count of hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashLen(key string) (length int, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		hash, err := s.getHash(tx, key)
		if err != nil {
			return err
		}
		length = len(hash)
		return nil
	})
	return
}

// HashKeys returns list of all hash fields. Error will occur if key doesn't exist or key type is not hash.
func (s *storage) HashKeys(key string) (keys []string, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		hash, err := s.getHash(tx, key)
		if err != nil {
			return err
		}
		for key := range hash {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return nil
	})
	return
}

// ListCreate creates new list with specified key and ttl. Use zero duration if key should exist forever.
func (s *storage) ListCreate(key string, ttl uint64) error {
	return notSupportedError
}

// ListLeftPop pops value from the list beginning.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListLeftPop(key string) (value string, err error) {
	return "", notSupportedError
}

// ListRightPop pops value from the list ending.
// Error will occur if key doesn't exist, key type is not list or list is empty.
func (s *storage) ListRightPop(key string) (value string, err error) {
	return "", notSupportedError
}

// ListLeftPush adds value to the list beginning. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLeftPush(key, value string) error {
	return notSupportedError
}

// ListRightPush adds value to the list ending. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListRightPush(key, value string) error {
	return notSupportedError
}

// ListLen returns count of elements in the list. Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListLen(key string) (length int, err error) {
	return 0, notSupportedError
}

// ListRange returns list of elements from the list from start to stop index.
// Error will occur if key doesn't exist or key type is not list.
func (s *storage) ListRange(key string, start, stop int) (values []string, err error) {
	return nil, notSupportedError
}
