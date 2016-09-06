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

func (s *StorageTestSuite) TestGetExpired(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.Set("key", "value", 1)
	time.Sleep(time.Second)

	_, err := storage.Get("key")
	c.Assert(err, ErrorMatches, "Key does not exist")
}

func (s *StorageTestSuite) TestLRU(c *C) {
	storage, _ := NewStorage(2, time.Minute)

	storage.Set("key1", "value1", 0)
	storage.Set("key2", "value2", 0)
	storage.Set("key3", "value3", 0)
	_, err := storage.Get("key1")
	c.Assert(err, ErrorMatches, "Key does not exist")
}

func (s *StorageTestSuite) TestKeys(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.Set("key1", "value1", 0)
	storage.Set("key2", "value2", 0)

	c.Assert(storage.Keys(), DeepEquals, []string{"key1", "key2"})
}

func (s *StorageTestSuite) TestSetAndGet(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Get non-existing key value and get error
	value1, err1 := storage.Get("key")
	c.Assert(err1, ErrorMatches, "Key does not exist")
	c.Assert(value1, Equals, "")

	// Set key value
	err2 := storage.Set("key", "value", 0)
	c.Assert(err2, IsNil)

	// Get existing key value
	value3, err3 := storage.Get("key")
	c.Assert(err3, IsNil)
	c.Assert(value3, Equals, "value")

	// Get string like hash and get error
	_, err4 := storage.HashGet("key", "field")
	c.Assert(err4, ErrorMatches, "Key type is not hash")

	// Try to set existing key value
	err5 := storage.Set("key", "value", 0)
	c.Assert(err5, ErrorMatches, "Key already exists")
}

func (s *StorageTestSuite) TestUpdate(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Update non-existing key value and get error
	err1 := storage.Update("key", "value")
	c.Assert(err1, ErrorMatches, "Key does not exist")

	// Set key value
	err2 := storage.Set("key", "value", 0)
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
	storage, _ := NewStorage(100, time.Minute)

	// Delete non-existing key and get error
	err1 := storage.Delete("key")
	c.Assert(err1, NotNil)

	// Set key value
	err2 := storage.Set("key", "value", 0)
	c.Assert(err2, IsNil)

	// Delete existing key
	err3 := storage.Delete("key")
	c.Assert(err3, IsNil)

	// Check that key doesn't exist anymore
	_, err4 := storage.Get("key")
	c.Assert(err4, NotNil)
}

func (s *StorageTestSuite) TestHashCreate(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Create hash
	err1 := storage.HashCreate("key1", 0)
	c.Assert(err1, IsNil)

	// Get hash like string and get error
	_, err2 := storage.Get("key1")
	c.Assert(err2, ErrorMatches, "Key type is not string")

	// Create hash with existing key
	storage.Set("key2", "value2", 0)
	err3 := storage.HashCreate("key2", 0)
	c.Assert(err3, ErrorMatches, "Key already exists")
}

func (s *StorageTestSuite) TestHashSetAndGet(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Get field from non-existing hash and get error
	value1, err1 := storage.HashGet("key", "field")
	c.Assert(err1, ErrorMatches, "Key does not exist")
	c.Assert(value1, Equals, "")

	// Create hash
	err2 := storage.HashCreate("key", 0)
	c.Assert(err2, IsNil)

	// Get non-existing field from existing hash
	value3, err3 := storage.HashGet("key", "field")
	c.Assert(err3, ErrorMatches, "Field does not exist")
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
	c.Assert(err6, IsNil)
}

func (s *StorageTestSuite) TestHashGetAll(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Create and fill hash
	storage.HashCreate("key", 0)
	storage.HashSet("key", "field", "value")

	hash, err1 := storage.HashGetAll("key")
	c.Assert(err1, IsNil)
	c.Assert(hash, DeepEquals, map[string]string{"field": "value"})

	_, err2 := storage.HashGetAll("key2")
	c.Assert(err2, ErrorMatches, "Key does not exist")
}

func (s *StorageTestSuite) TestHashDelete(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Create hash
	err1 := storage.HashCreate("key", 0)
	c.Assert(err1, IsNil)

	// Delete non-existing field and get error
	err2 := storage.HashDelete("key", "field")
	c.Assert(err2, ErrorMatches, "Field does not exist")

	// Set hash field
	err3 := storage.HashSet("key", "field", "value")
	c.Assert(err3, IsNil)

	// Delete existing field
	err4 := storage.HashDelete("key", "field")
	c.Assert(err4, IsNil)

	// Chech field is deleted
	_, err5 := storage.HashGet("key", "field")
	c.Assert(err5, ErrorMatches, "Field does not exist")

	// Delete non-existing hash field
	err6 := storage.HashDelete("key2", "field")
	c.Assert(err6, ErrorMatches, "Key does not exist")
}

func (s *StorageTestSuite) TestHashLenAndKeys(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Create and fill hash
	storage.HashCreate("key", 0)
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
	c.Assert(err3, ErrorMatches, "Key does not exist")
	_, err4 := storage.HashKeys("key2")
	c.Assert(err4, ErrorMatches, "Key does not exist")

	// Get length and keys of non-hash value
	storage.Set("key3", "value3", 0)
	_, err5 := storage.HashLen("key3")
	c.Assert(err5, ErrorMatches, "Key type is not hash")
	_, err6 := storage.HashKeys("key3")
	c.Assert(err6, ErrorMatches, "Key type is not hash")
}

func (s *StorageTestSuite) TestListCreate(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Create list
	err1 := storage.ListCreate("key1", 0)
	c.Assert(err1, IsNil)

	// Get list like string and get error
	_, err2 := storage.Get("key1")
	c.Assert(err2, ErrorMatches, "Key type is not string")

	// Create list with existing key
	storage.Set("key2", "value2", 0)
	err3 := storage.ListCreate("key2", 0)
	c.Assert(err3, ErrorMatches, "Key already exists")
}

func (s *StorageTestSuite) TestListLen(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	// Create list
	storage.ListCreate("key", 0)
	storage.ListRightPush("key", "value1")
	storage.ListRightPush("key", "value2")

	len1, err1 := storage.ListLen("key")
	c.Assert(err1, IsNil)
	c.Assert(len1, Equals, 2)

	// Get length of non-existing list and non-list value
	storage.Set("key2", "value2", 0)
	_, err2 := storage.ListLen("key2")
	c.Assert(err2, ErrorMatches, "Key type is not list")
	_, err3 := storage.ListLen("key3")
	c.Assert(err3, ErrorMatches, "Key does not exist")
}

func (s *StorageTestSuite) TestListPushAndPop(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.ListCreate("key", 0)

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
	c.Assert(err7, ErrorMatches, "List is empty")
	_, err8 := storage.ListLeftPop("key")
	c.Assert(err8, ErrorMatches, "List is empty")

	// Check push and pop operation on non-existing list and non-list key
	var err error
	storage.Set("key2", "value2", 0)
	_, err = storage.ListLeftPop("key2")
	c.Assert(err, ErrorMatches, "Key type is not list")
	_, err = storage.ListLeftPop("key3")
	c.Assert(err, ErrorMatches, "Key does not exist")
	_, err = storage.ListRightPop("key2")
	c.Assert(err, ErrorMatches, "Key type is not list")
	_, err = storage.ListRightPop("key3")
	c.Assert(err, ErrorMatches, "Key does not exist")
	err = storage.ListLeftPush("key2", "value2")
	c.Assert(err, ErrorMatches, "Key type is not list")
	err = storage.ListLeftPush("key3", "value3")
	c.Assert(err, IsNil)
	err = storage.ListRightPush("key2", "value2")
	c.Assert(err, ErrorMatches, "Key type is not list")
	err = storage.ListRightPush("key3", "value3")
	c.Assert(err, IsNil)
}

func (s *StorageTestSuite) TestListRange(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.ListCreate("key", 0)
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
	storage.Set("key2", "value2", 0)
	_, err = storage.ListRange("key2", 0, 0)
	c.Assert(err, ErrorMatches, "Key type is not list")
	_, err = storage.ListRange("key3", 0, 0)
	c.Assert(err, ErrorMatches, "Key does not exist")
}

func (s *StorageTestSuite) TestGC(c *C) {
	storage, _ := NewStorage(100, time.Millisecond)
	storage.Set("key", "value", 1)
	c.Assert(storage.lru.Len(), Equals, 1)
	// Wait for expiring + gc interval
	time.Sleep(time.Second + time.Millisecond)
	c.Assert(storage.lru.Len(), Equals, 0)
}

func (s *StorageTestSuite) BenchmarkGet(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.Set("key", "value", 0)

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		storage.Get("key")
	}
}

func (s *StorageTestSuite) BenchmarkGetNotExist(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		storage.Get("key")
	}
}

func (s *StorageTestSuite) BenchmarkSetExist(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.Set("key", "value", 0)

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		storage.Set("key", "value", 0)
	}
}

func (s *StorageTestSuite) BenchmarkUpdate(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.Set("key", "value", 0)

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		storage.Update("key", "value")
	}
}

func (s *StorageTestSuite) BenchmarkHashGetAll(c *C) {
	storage, _ := NewStorage(100, time.Minute)

	storage.HashCreate("key", 0)

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		storage.HashGetAll("key")
	}
}
