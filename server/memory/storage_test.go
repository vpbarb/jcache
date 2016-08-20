package memory

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
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
	c.Assert(err3, ErrorMatches, `Key "key3" does not exist`)
}

func (s *StorageTestSuite) TestSetAndGet(c *C) {
	storage := NewStorage()

	// Get non-existing key value and get error
	value1, err1 := storage.Get("key")
	c.Assert(err1, ErrorMatches, `Key "key" does not exist`)
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

func (s *StorageTestSuite) TestUpdate(c *C) {
	storage := NewStorage()

	// Update non-existing key value and get error
	err1 := storage.Update("key", "value")
	c.Assert(err1, ErrorMatches, `Key "key" does not exist`)

	// Set key value
	err2 := storage.Set("key", "value", time.Duration(0))
	c.Assert(err2, IsNil)

	// Update existing key value
	err3 := storage.Update("key", "updated")
	c.Assert(err3, IsNil)

	// Check updated value
	value4, err4 := storage.Get("key")
	c.Assert(err4, IsNil)
	c.Assert(value4, Equals, "updated")
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

func (s *StorageTestSuite) TestListCreate(c *C) {
	storage := NewStorage()

	// Create list
	err1 := storage.ListCreate("key1", time.Duration(0))
	c.Assert(err1, IsNil)

	// Get list like string and get error
	_, err2 := storage.Get("key1")
	c.Assert(err2, ErrorMatches, "Key type is not string")

	// Create list with existing key
	storage.Set("key2", "value2", time.Duration(0))
	err3 := storage.ListCreate("key2", time.Duration(0))
	c.Assert(err3, ErrorMatches, `Key "key2" already exists`)
}

func (s *StorageTestSuite) TestListLen(c *C) {
	storage := NewStorage()

	// Create list
	storage.ListCreate("key", time.Duration(0))
	storage.ListRightPush("key", "value1")
	storage.ListRightPush("key", "value2")

	len1, err1 := storage.ListLen("key")
	c.Assert(err1, IsNil)
	c.Assert(len1, Equals, 2)

	// Get length of non-existing list and non-list value
	storage.Set("key2", "value2", time.Duration(0))
	_, err2 := storage.ListLen("key2")
	c.Assert(err2, ErrorMatches, `Key type is not list`)
	_, err3 := storage.ListLen("key3")
	c.Assert(err3, ErrorMatches, `Key "key3" does not exist`)
}

func (s *StorageTestSuite) TestListPushAndPop(c *C) {
	storage := NewStorage()

	storage.ListCreate("key", time.Duration(0))

	// Make "1"
	err1 := storage.ListLeftPush("key", "1")
	c.Assert(err1, IsNil)

	// Make "1, 2"
	err2 := storage.ListRightPush("key", "2")
	c.Assert(err2, IsNil)

	// Make "3, 1, 2"
	err3 := storage.ListLeftPush("key", "3")
	c.Assert(err3, IsNil)

	// Make "3, 1"
	value4, err4 := storage.ListRightPop("key")
	c.Assert(err4, IsNil)
	c.Assert(value4, Equals, "2")

	// Make "1"
	value5, err5 := storage.ListLeftPop("key")
	c.Assert(err5, IsNil)
	c.Assert(value5, Equals, "3")

	// Make empty
	value6, err6 := storage.ListRightPop("key")
	c.Assert(err6, IsNil)
	c.Assert(value6, Equals, "1")

	// Check error on empty list
	_, err7 := storage.ListRightPop("key")
	c.Assert(err7, ErrorMatches, `List "key" is empty`)
	_, err8 := storage.ListLeftPop("key")
	c.Assert(err8, ErrorMatches, `List "key" is empty`)

	// Check push and pop operation on non-existing list and non-list key
	var err error
	storage.Set("key2", "value2", time.Duration(0))
	_, err = storage.ListLeftPop("key2")
	c.Assert(err, ErrorMatches, `Key type is not list`)
	_, err = storage.ListLeftPop("key3")
	c.Assert(err, ErrorMatches, `Key "key3" does not exist`)
	_, err = storage.ListRightPop("key2")
	c.Assert(err, ErrorMatches, `Key type is not list`)
	_, err = storage.ListRightPop("key3")
	c.Assert(err, ErrorMatches, `Key "key3" does not exist`)
	err = storage.ListLeftPush("key2", "value2")
	c.Assert(err, ErrorMatches, `Key type is not list`)
	err = storage.ListLeftPush("key3", "value3")
	c.Assert(err, ErrorMatches, `Key "key3" does not exist`)
	err = storage.ListRightPush("key2", "value2")
	c.Assert(err, ErrorMatches, `Key type is not list`)
	err = storage.ListRightPush("key3", "value3")
	c.Assert(err, ErrorMatches, `Key "key3" does not exist`)
}

func (s *StorageTestSuite) TestListRange(c *C) {
	storage := NewStorage()

	storage.ListCreate("key", time.Duration(0))
	storage.ListRightPush("key", "0")
	storage.ListRightPush("key", "1")
	storage.ListRightPush("key", "2")
	storage.ListRightPush("key", "3")
	storage.ListRightPush("key", "4")

	values1, err1 := storage.ListRange("key", 0, 0)
	c.Assert(err1, IsNil)
	c.Assert(values1, DeepEquals, []string{"0"})

	values2, err2 := storage.ListRange("key", 1, 2)
	c.Assert(err2, IsNil)
	c.Assert(values2, DeepEquals, []string{"1", "2"})

	values3, err3 := storage.ListRange("key", -1, 1)
	c.Assert(err3, IsNil)
	c.Assert(values3, DeepEquals, []string{"0", "1"})

	values4, err4 := storage.ListRange("key", 3, 100)
	c.Assert(err4, IsNil)
	c.Assert(values4, DeepEquals, []string{"3", "4"})

	// Get range of non-existing list and non-list value
	var err error
	storage.Set("key2", "value2", time.Duration(0))
	_, err = storage.ListRange("key2", 0, 0)
	c.Assert(err, ErrorMatches, `Key type is not list`)
	_, err = storage.ListRange("key3", 0, 0)
	c.Assert(err, ErrorMatches, `Key "key3" does not exist`)
}
