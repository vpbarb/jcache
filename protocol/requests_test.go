package protocol

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type RequestsTestSuite struct{}

var _ = Suite(&RequestsTestSuite{})

func (s *RequestsTestSuite) TestRequestEncode(c *C) {
	request := newRequest("CMD")
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD\r\n"))
}

func (s *RequestsTestSuite) TestRequestDecode(c *C) {
	request := newRequest("CMD")
	err := request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, IsNil)
}

func (s *RequestsTestSuite) TestRequestDecodeError(c *C) {
	request := newRequest("CMD")
	err := request.Decode([]byte("TEST"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestKeyEncode(c *C) {
	request := newKeyRequest("CMD")
	request.Key = "key"
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD key\r\n"))
}

func (s *RequestsTestSuite) TestKeyEncodeError(c *C) {
	request := newKeyRequest("CMD")
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")
}

func (s *RequestsTestSuite) TestKeyDecode(c *C) {
	request := newKeyRequest("CMD")
	err := request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, IsNil)
	c.Assert(request.Key, Equals, "key")
}

func (s *RequestsTestSuite) TestKeyDecodeError(c *C) {
	request := newKeyRequest("CMD")
	var err error
	err = request.Decode([]byte("TEST key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD +"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestSetEncode(c *C) {
	request := &setRequest{keyValueRequest: newKeyValueRequest("CMD")}
	request.Key = "key"
	request.Value = "value"
	request.TTL = 3
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD key 3 5\r\nvalue\r\n"))
}

func (s *RequestsTestSuite) TestSetEncodeError(c *C) {
	request := &setRequest{keyValueRequest: newKeyValueRequest("CMD")}
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")
}

func (s *RequestsTestSuite) TestSetDecode(c *C) {
	request := &setRequest{keyValueRequest: newKeyValueRequest("CMD")}
	err := request.Decode([]byte("CMD key 3 5"), bytes.NewBufferString("value\r\n"))
	c.Assert(err, IsNil)
	c.Assert(request.Key, Equals, "key")
	c.Assert(request.Value, Equals, "value")
	c.Assert(request.TTL, Equals, uint64(3))
}

func (s *RequestsTestSuite) TestSetDecodeError(c *C) {
	request := &setRequest{keyValueRequest: newKeyValueRequest("CMD")}
	var err error
	err = request.Decode([]byte("TEST key 3 5"), bytes.NewBufferString("value\r\n"))
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key 0"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key 0 str"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key 3 5"), bytes.NewBufferString("v\r\n"))
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestKeyValueEncode(c *C) {
	request := newKeyValueRequest("CMD")
	request.Key = "key"
	request.Value = "value"
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD key 5\r\nvalue\r\n"))
}

func (s *RequestsTestSuite) TestKeyValueEncodeError(c *C) {
	request := newKeyValueRequest("CMD")
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")
}

func (s *RequestsTestSuite) TestKeyValueDecode(c *C) {
	request := newKeyValueRequest("CMD")
	err := request.Decode([]byte("CMD key 5"), bytes.NewBufferString("value\r\n"))
	c.Assert(err, IsNil)
	c.Assert(request.Key, Equals, "key")
	c.Assert(request.Value, Equals, "value")
}

func (s *RequestsTestSuite) TestKeyValueDecodeError(c *C) {
	request := newKeyValueRequest("CMD")
	var err error
	err = request.Decode([]byte("TEST key 5"), bytes.NewBufferString("value\r\n"))
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key str"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key 5"), bytes.NewBufferString("v\r\n"))
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestKeyFieldEncode(c *C) {
	request := newKeyFieldRequest("CMD")
	request.Key = "key"
	request.Field = "field"
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD key field\r\n"))
}

func (s *RequestsTestSuite) TestKeyFieldEncodeError(c *C) {
	request := newKeyFieldRequest("CMD")
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "key"
	request.Field = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Field is not valid")
}

func (s *RequestsTestSuite) TestKeyFieldDecode(c *C) {
	request := newKeyFieldRequest("CMD")
	err := request.Decode([]byte("CMD key field"), &bytes.Buffer{})
	c.Assert(err, IsNil)
	c.Assert(request.Key, Equals, "key")
	c.Assert(request.Field, Equals, "field")
}

func (s *RequestsTestSuite) TestKeyFieldDecodeError(c *C) {
	request := newKeyFieldRequest("CMD")
	var err error
	err = request.Decode([]byte("TEST key field"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key +"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestKeyTTLEncode(c *C) {
	request := newKeyTTLRequest("CMD")
	request.Key = "key"
	request.TTL = 5
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD key 5\r\n"))
}

func (s *RequestsTestSuite) TestKeyTTLEncodeError(c *C) {
	request := newKeyTTLRequest("CMD")
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")
}

func (s *RequestsTestSuite) TestKeyTTLDecode(c *C) {
	request := newKeyTTLRequest("CMD")
	err := request.Decode([]byte("CMD key 5"), &bytes.Buffer{})
	c.Assert(err, IsNil)
	c.Assert(request.Key, Equals, "key")
	c.Assert(request.TTL, Equals, uint64(5))
}

func (s *RequestsTestSuite) TestKeyTTLDecodeError(c *C) {
	request := newKeyTTLRequest("CMD")
	var err error
	err = request.Decode([]byte("TEST key 5"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key str"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestKeyFieldValueEncode(c *C) {
	request := newKeyFieldValueRequest("CMD")
	request.Key = "key"
	request.Field = "field"
	request.Value = "value"
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD key field 5\r\nvalue\r\n"))
}

func (s *RequestsTestSuite) TestKeyFieldValueEncodeError(c *C) {
	request := newKeyFieldValueRequest("CMD")
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "key"
	request.Field = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Field is not valid")
}

func (s *RequestsTestSuite) TestKeyFieldValueDecode(c *C) {
	request := newKeyFieldValueRequest("CMD")
	err := request.Decode([]byte("CMD key field 12"), bytes.NewBufferString("value\r\nvalue\r\n"))
	c.Assert(err, IsNil)
	c.Assert(request.Key, Equals, "key")
	c.Assert(request.Field, Equals, "field")
	c.Assert(request.Value, Equals, "value\r\nvalue")
}

func (s *RequestsTestSuite) TestKeyFieldValueDecodeError(c *C) {
	request := newKeyFieldValueRequest("CMD")
	var err error
	err = request.Decode([]byte("TEST key field 12"), bytes.NewBufferString("value\r\nvalue\r\n"))
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key +"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key field"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key field str"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key field 20"), bytes.NewBufferString("value\r\nvalue\r\n"))
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestListRangeEncode(c *C) {
	request := &listRangeRequest{keyRequest: newKeyRequest("CMD")}
	request.Key = "key"
	request.Start = 1
	request.Stop = 3
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD key 1 3\r\n"))
}

func (s *RequestsTestSuite) TestListRangeEncodeError(c *C) {
	request := &listRangeRequest{keyRequest: newKeyRequest("CMD")}
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")

	request.Key = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Key is not valid")
}

func (s *RequestsTestSuite) TestListRangeDecode(c *C) {
	request := &listRangeRequest{keyRequest: newKeyRequest("CMD")}
	err := request.Decode([]byte("CMD key 1 3"), &bytes.Buffer{})
	c.Assert(err, IsNil)
	c.Assert(request.Key, Equals, "key")
	c.Assert(request.Start, Equals, 1)
	c.Assert(request.Stop, Equals, 3)
}

func (s *RequestsTestSuite) TestListRangeDecodeError(c *C) {
	request := &listRangeRequest{keyRequest: newKeyRequest("CMD")}
	var err error
	err = request.Decode([]byte("TEST key 0 0"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key 1"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD key 1 str"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
}

func (s *RequestsTestSuite) TestAuthEncode(c *C) {
	request := &authRequest{request: newRequest("CMD")}
	request.User = "user"
	request.Password = "password"
	data, err := request.Encode()
	c.Assert(err, IsNil)
	c.Assert(data, DeepEquals, []byte("CMD user password\r\n"))
}

func (s *RequestsTestSuite) TestAuthEncodeError(c *C) {
	request := &authRequest{request: newRequest("CMD")}
	_, err := request.Encode()
	c.Assert(err, ErrorMatches, "User is not valid")

	request.User = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "User is not valid")

	request.User = "user"
	request.Password = "+"
	_, err = request.Encode()
	c.Assert(err, ErrorMatches, "Password is not valid")
}

func (s *RequestsTestSuite) TestAuthDecode(c *C) {
	request := &authRequest{request: newRequest("CMD")}
	err := request.Decode([]byte("CMD user password"), &bytes.Buffer{})
	c.Assert(err, IsNil)
	c.Assert(request.User, Equals, "user")
	c.Assert(request.Password, Equals, "password")
}

func (s *RequestsTestSuite) TestAuthDecodeError(c *C) {
	request := &authRequest{request: newRequest("CMD")}
	var err error
	err = request.Decode([]byte("TEST user password"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command")
	err = request.Decode([]byte("CMD"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD user"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
	err = request.Decode([]byte("CMD user +"), &bytes.Buffer{})
	c.Assert(err, ErrorMatches, "Invalid command format")
}