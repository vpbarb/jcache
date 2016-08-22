package storage

import (
	"container/list"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ItemTestSuite struct{}

var _ = Suite(&ItemTestSuite{})

func (s *ItemTestSuite) TestStringWithTTL(c *C) {
	item := NewItem("value", 10)

	value, err := item.CastString()
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "value")

	c.Assert(item.IsAlive(), Equals, true)

	_, err = item.CastHash()
	c.Assert(err, NotNil)
	_, err = item.CastList()
	c.Assert(err, NotNil)
}

func (s *ItemTestSuite) TestExpired(c *C) {
	item := NewItem("value", 1)

	time.Sleep(time.Second)

	c.Assert(item.IsAlive(), Equals, false)
}

func (s *ItemTestSuite) TestHash(c *C) {
	item := NewItem(Hash{"field": "value"}, 0)

	hash, err := item.CastHash()
	c.Assert(err, IsNil)
	c.Assert(hash, DeepEquals, Hash{"field": "value"})

	c.Assert(item.IsAlive(), Equals, true)

	_, err = item.CastString()
	c.Assert(err, NotNil)
	_, err = item.CastList()
	c.Assert(err, NotNil)
}

func (s *ItemTestSuite) TestList(c *C) {
	list := &list.List{}
	item := NewItem(list, 0)

	list, err := item.CastList()
	c.Assert(err, IsNil)
	c.Assert(list, FitsTypeOf, list)

	c.Assert(item.IsAlive(), Equals, true)

	_, err = item.CastString()
	c.Assert(err, NotNil)
	_, err = item.CastHash()
	c.Assert(err, NotNil)
}
