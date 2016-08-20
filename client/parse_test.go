package client

import (
	. "gopkg.in/check.v1"
)

type ParseTestSuite struct{}

var _ = Suite(&ParseTestSuite{})

func (s *ParseTestSuite) TestParseValue(c *C) {
	value, err := parseValue(`"some value"`)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, "some value")
}

func (s *ParseTestSuite) TestParseValueError(c *C) {
	value, err := parseValue(`"some value without quote`)
	c.Assert(err, NotNil)
	c.Assert(value, Equals, "")
}

func (s *ParseTestSuite) TestParseHashField(c *C) {
	field, value, err := parseHashField(`some_field:"some value"`)
	c.Assert(err, IsNil)
	c.Assert(field, Equals, "some_field")
	c.Assert(value, Equals, "some value")
}

func (s *ParseTestSuite) TestParseHashFieldError(c *C) {
	field, value, err := parseHashField(`some_field:"some value without quote`)
	c.Assert(err, NotNil)
	c.Assert(field, Equals, "")
	c.Assert(value, Equals, "")
}
