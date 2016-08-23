package protocol

import (
	//"bytes"
	"errors"

	"bytes"
	. "gopkg.in/check.v1"
)

type ResponsesTestSuite struct{}

var _ = Suite(&ResponsesTestSuite{})

func (s *ResponsesTestSuite) TestOkEncode(c *C) {
	response := newOkResponse()

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("OK\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestOkDecode(c *C) {
	response := newOkResponse()

	err := response.Decode(bytes.NewBufferString("OK\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestOkDecodeError(c *C) {
	response := newOkResponse()

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
}

func (s *ResponsesTestSuite) TestLenEncode(c *C) {
	response := &lenResponse{dataResponse: newDataResponse()}
	response.Len = 5

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("DATA\r\nLEN 5\r\nEND\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestLenDecode(c *C) {
	response := &lenResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("DATA\r\nLEN 5\r\nEND\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)
	c.Assert(response.Len, Equals, 5)

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestLenDecodeError(c *C) {
	response := &lenResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("DATA\r\nLEN abc\r\nEND\r\n"))
	c.Assert(err, ErrorMatches, "Invalid data format")
}

func (s *ResponsesTestSuite) TestValueEncode(c *C) {
	response := &valueResponse{dataResponse: newDataResponse()}
	response.Value = "value"

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("DATA\r\nVALUE 5\r\nvalue\r\nEND\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestValueDecode(c *C) {
	response := &valueResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("DATA\r\nVALUE 5\r\nvalue\r\nEND\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)
	c.Assert(response.Value, Equals, "value")

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestValueDecodeError(c *C) {
	response := &valueResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("DATA\r\nVALUE 5\r\nval\r\nEND\r\n"))
	c.Assert(err, ErrorMatches, "Value length is invalid")
}
