package memory

import (
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

type StorageTestSuite struct{}

var _ = Suite(&StorageTestSuite{})

func Test(t *testing.T) {
	TestingT(t)
}

func (s *StorageTestSuite) TestExpire(c *C) {
	storage := NewStorage()

	storage.Set("key", "value", time.Millisecond)
	time.Sleep(2 * time.Millisecond)

	_, err := storage.Get("key")
	c.Assert(err, ErrorMatches, `Key "key" does not exist`)
}

func (s *StorageTestSuite) TestKeys(c *C) {
	storage := NewStorage()

	storage.Set("key1", "value1", time.Duration(0))
	storage.Set("key2", "value2", time.Duration(0))

	c.Assert(storage.Keys(), DeepEquals, []string{"key1", "key2"})
}

func (s *StorageTestSuite) TestTTL(c *C) {
	storage := NewStorage()

	storage.Set("key1", "value1", time.Hour)
	storage.Set("key2", "value2", time.Duration(0))

	ttl1, err1 := storage.TTL("key1")
	c.Assert(err1, IsNil)
	c.Assert(ttl1 > time.Duration(0), Equals, true)

	ttl2, err2 := storage.TTL("key2")
	c.Assert(err2, IsNil)
	c.Assert(ttl2 == time.Duration(0), Equals, true)

	_, err3 := storage.TTL("key3")
	c.Assert(err3, NotNil)
}

func (s *StorageTestSuite) TestSetAndGet(c *C) {
	storage := NewStorage()

	// Get non-existing key value and get error
	value1, err1 := storage.Get("key")
	c.Assert(err1, NotNil)
	c.Assert(value1, Equals, "")

	// Set key value
	err2 := storage.Set("key", "value", time.Duration(0))
	c.Assert(err2, IsNil)

	// Get existing key value
	value3, err3 := storage.Get("key")
	c.Assert(err3, IsNil)
	c.Assert(value3, Equals, "value")

	// Get string like hash and get error
	_, err4 := storage.HashGet("key", "field")
	c.Assert(err4, ErrorMatches, "Key type is not hash")

	// Try to set existing key value
	err5 := storage.Set("key", "value", time.Duration(0))
	c.Assert(err5, ErrorMatches, `Key "key" already exists`)
}

func (s *StorageTestSuite) TestDelete(c *C) {
	storage := NewStorage()

	// Delete non-existing key and get error
	err1 := storage.Delete("key")
	c.Assert(err1, NotNil)

	// Set key value
	err2 := storage.Set("key", "value", time.Duration(0))
	c.Assert(err2, IsNil)

	// Delete existing key
	err3 := storage.Delete("key")
	c.Assert(err3, IsNil)

	// Check that key doesn't exist anymore
	_, err4 := storage.Get("key")
	c.Assert(err4, NotNil)
}

func (s *StorageTestSuite) TestHashCreate(c *C) {
	storage := NewStorage()

	// Create hash
	err1 := storage.HashCreate("key1", time.Duration(0))
	c.Assert(err1, IsNil)

	// Get hash like string and get error
	_, err2 := storage.Get("key1")
	c.Assert(err2, ErrorMatches, "Key type is not string")

	// Create hash with existing key
	storage.Set("key2", "value2", time.Duration(0))
	err3 := storage.HashCreate("key2", time.Duration(0))
	c.Assert(err3, ErrorMatches, `Key "key2" already exists`)
}

func (s *StorageTestSuite) TestHashSetAndGet(c *C) {
	storage := NewStorage()

	// Get field from non-existing hash and get error
	value1, err1 := storage.HashGet("key", "field")
	c.Assert(err1, ErrorMatches, `Key "key" does not exist`)
	c.Assert(value1, Equals, "")

	// Create hash
	err2 := storage.HashCreate("key", time.Duration(0))
	c.Assert(err2, IsNil)

	// Get non-existing field from existing hash
	value3, err3 := storage.HashGet("key", "field")
	c.Assert(err3, ErrorMatches, `Field "field" does not exist`)
	c.Assert(value3, Equals, "")

	// Set hash field
	err4 := storage.HashSet("key", "field", "value")
	c.Assert(err4, IsNil)

	// Get existing field from existing hash
	value5, err5 := storage.HashGet("key", "field")
	c.Assert(err5, IsNil)
	c.Assert(value5, Equals, "value")

	// Set non-existing hash field
	err6 := storage.HashSet("key2", "field", "value")
	c.Assert(err6, ErrorMatches, `Key "key2" does not exist`)
}

func (s *StorageTestSuite) TestHashGetAll(c *C) {
	storage := NewStorage()

	// Create and fill hash
	storage.HashCreate("key", time.Duration(0))
	storage.HashSet("key", "field", "value")

	hash, err1 := storage.HashGetAll("key")
	c.Assert(err1, IsNil)
	c.Assert(hash, DeepEquals, map[string]string{"field": "value"})

	_, err2 := storage.HashGetAll("key2")
	c.Assert(err2, ErrorMatches, `Key "key2" does not exist`)
}

func (s *StorageTestSuite) TestHashDelete(c *C) {
	storage := NewStorage()

	// Create hash
	err1 := storage.HashCreate("key", time.Duration(0))
	c.Assert(err1, IsNil)

	// Delete non-existing field and get error
	err2 := storage.HashDelete("key", "field")
	c.Assert(err2, ErrorMatches, `Field "field" does not exist`)

	// Set hash field
	err3 := storage.HashSet("key", "field", "value")
	c.Assert(err3, IsNil)

	// Delete existing field
	err4 := storage.HashDelete("key", "field")
	c.Assert(err4, IsNil)

	// Chech field is deleted
	_, err5 := storage.HashGet("key", "field")
	c.Assert(err5, ErrorMatches, `Field "field" does not exist`)

	// Delete non-existing hash field
	err6 := storage.HashDelete("key2", "field")
	c.Assert(err6, ErrorMatches, `Key "key2" does not exist`)
}

func (s *StorageTestSuite) TestHashLenAndKeys(c *C) {
	storage := NewStorage()

	// Create and fill hash
	storage.HashCreate("key", time.Duration(0))
	storage.HashSet("key", "field1", "value1")
	storage.HashSet("key", "field2", "value2")

	len, err1 := storage.HashLen("key")
	c.Assert(err1, IsNil)
	c.Assert(len, Equals, 2)

	keys, err2 := storage.HashKeys("key")
	c.Assert(err2, IsNil)
	c.Assert(keys, DeepEquals, []string{"field1", "field2"})

	// Get length and keys of non-existing hash
	_, err3 := storage.HashLen("key2")
	c.Assert(err3, ErrorMatches, `Key "key2" does not exist`)
	_, err4 := storage.HashKeys("key2")
	c.Assert(err4, ErrorMatches, `Key "key2" does not exist`)

	// Get length and keys of non-hash value
	storage.Set("key3", "value3", time.Duration(0))
	_, err5 := storage.HashLen("key3")
	c.Assert(err5, ErrorMatches, `Key type is not hash`)
	_, err6 := storage.HashKeys("key3")
	c.Assert(err6, ErrorMatches, `Key type is not hash`)
}
