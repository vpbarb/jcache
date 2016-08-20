package multi

import (
	"time"

	"github.com/Barberrrry/jcache/server/storage/memory"
	. "gopkg.in/check.v1"
)

type MultiStorageTestSuite struct{}

var _ = Suite(&MultiStorageTestSuite{})

func (s *MultiStorageTestSuite) TestKeys(c *C) {
	storage := NewStorage()
	storage.AddStorage(memory.NewStorage())

	storage.Set("key0", "value0", time.Duration(0))
	storage.Set("key1", "value1", time.Duration(0))
	storage.Set("key2", "value2", time.Duration(0))
	storage.Set("key3", "value3", time.Duration(0))
	storage.Set("key4", "value4", time.Duration(0))
	storage.Set("key5", "value5", time.Duration(0))
	storage.Set("key6", "value6", time.Duration(0))
	storage.Set("key7", "value7", time.Duration(0))
	storage.Set("key8", "value8", time.Duration(0))
	storage.Set("key9", "value9", time.Duration(0))

	c.Assert(storage.Keys(), DeepEquals, []string{
		"key0",
		"key1",
		"key2",
		"key3",
		"key4",
		"key5",
		"key6",
		"key7",
		"key8",
		"key9",
	})
}

func (s *MultiStorageTestSuite) TestSetAndGet(c *C) {
	storage := NewStorage()
	storage.AddStorage(memory.NewStorage())

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
