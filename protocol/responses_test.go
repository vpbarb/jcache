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
	response := &lenResponse{response: &response{}}
	response.Len = 5

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("LEN 5\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestLenDecode(c *C) {
	response := &lenResponse{response: &response{}}

	err := response.Decode(bytes.NewBufferString("LEN 5\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)
	c.Assert(response.Len, Equals, 5)

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestLenDecodeError(c *C) {
	response := &lenResponse{response: &response{}}

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("LEN abc\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
}

func (s *ResponsesTestSuite) TestValueEncode(c *C) {
	response := newValueResponse()
	response.Value = "value"

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("VALUE 5\r\nvalue\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestValueDecode(c *C) {
	response := newValueResponse()

	err := response.Decode(bytes.NewBufferString("VALUE 5\r\nvalue\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)
	c.Assert(response.Value, Equals, "value")

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestValueDecodeError(c *C) {
	response := newValueResponse()

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("VALUE 5\r\nval\r\n"))
	c.Assert(err, ErrorMatches, "Value length is invalid")
}

func (s *ResponsesTestSuite) TestKeysEncode(c *C) {
	response := &keysResponse{dataResponse: newDataResponse()}
	response.Keys = []string{"key1", "key2"}

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("DATA\r\nKEY key1\r\nKEY key2\r\nEND\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestKeysDecode(c *C) {
	response := &keysResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("DATA\r\nKEY key1\r\nKEY key2\r\nEND\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)
	c.Assert(response.Keys, DeepEquals, []string{"key1", "key2"})

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestKeysDecodeError(c *C) {
	response := &keysResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("DATA\r\nVALUE key\r\nEND\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
}

func (s *ResponsesTestSuite) TestFieldsEncode(c *C) {
	response := &fieldsResponse{dataResponse: newDataResponse()}
	response.Fields = map[string]string{"field1": "value1", "field2": "value2"}

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("DATA\r\nFIELD field1 6\r\nvalue1\r\nFIELD field2 6\r\nvalue2\r\nEND\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestFieldsDecode(c *C) {
	response := &fieldsResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("DATA\r\nFIELD field1 6\r\nvalue1\r\nFIELD field2 6\r\nvalue2\r\nEND\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)
	c.Assert(response.Fields["field1"], Equals, "value1")
	c.Assert(response.Fields["field2"], Equals, "value2")

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestFieldsDecodeError(c *C) {
	response := &fieldsResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("DATA\r\nFIELD field\r\nEND\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("DATA\r\nFIELD field 10\r\nEND\r\n"))
	c.Assert(err, ErrorMatches, "Value length is invalid")
}

func (s *ResponsesTestSuite) TestValuesEncode(c *C) {
	response := &valuesResponse{dataResponse: newDataResponse()}
	response.Values = []string{"value1", "value100"}

	data, err := response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("DATA\r\nVALUE 6\r\nvalue1\r\nVALUE 8\r\nvalue100\r\nEND\r\n"))

	response.Error = errors.New("TEST")
	data, err = response.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("ERROR TEST\r\n"))
}

func (s *ResponsesTestSuite) TestValuesDecode(c *C) {
	response := &valuesResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("DATA\r\nVALUE 6\r\nvalue1\r\nVALUE 8\r\nvalue100\r\nEND\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, IsNil)
	c.Assert(response.Values, DeepEquals, []string{"value1", "value100"})

	err = response.Decode(bytes.NewBufferString("ERROR TEST\r\n"))
	c.Assert(err, IsNil)
	c.Assert(response.Error, ErrorMatches, "Response error: TEST")
}

func (s *ResponsesTestSuite) TestValuesDecodeError(c *C) {
	response := &valuesResponse{dataResponse: newDataResponse()}

	err := response.Decode(bytes.NewBufferString("TEST\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("DATA\r\nVALUE key\r\nEND\r\n"))
	c.Assert(err, ErrorMatches, "Invalid response format")
	err = response.Decode(bytes.NewBufferString("DATA\r\nVALUE 5\r\nEND\r\n"))
	c.Assert(err, ErrorMatches, "Value length is invalid")
}
