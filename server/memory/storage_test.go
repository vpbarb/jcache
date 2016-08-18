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

func (s *StorageTestSuite) TestHashCreateAndSetAndGet(c *C) {
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
}
