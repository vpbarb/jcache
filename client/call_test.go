package client

import (
	"bytes"
	"errors"

	. "gopkg.in/check.v1"
)

type CallTestSuite struct{}

var _ = Suite(&CallTestSuite{})

type errorWriter struct {
	err error
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}

func (s *CallTestSuite) TestWriteError(c *C) {
	w := &errorWriter{errors.New("Write error")}
	r := &bytes.Buffer{}

	response, err := call(w, r, "TEST")
	c.Assert(err, ErrorMatches, "Cannot write to connection: Write error")
	c.Assert(response, IsNil)
}

func (s *CallTestSuite) TestWriteOk(c *C) {
	w := &bytes.Buffer{}
	r := &bytes.Buffer{}

	call(w, r, "TEST")
	c.Assert(w.String(), DeepEquals, "TEST\r\n")
}

func (s *CallTestSuite) TestReadError(c *C) {
	w := &bytes.Buffer{}
	r := &bytes.Buffer{}

	response, err := call(w, r, "TEST")
	c.Assert(err, ErrorMatches, "Cannot read from connection: EOF")
	c.Assert(response, IsNil)
}

//func (s *CallTestSuite) TestInvalidRowsCount(c *C) {
//	w := &bytes.Buffer{}
//	r := bytes.NewBufferString("$a\r\n")
//
//	response, err := call(w, r, "TEST")
//	c.Assert(err, ErrorMatches, "Invalid response rows count")
//	c.Assert(response, IsNil)
//}

func (s *CallTestSuite) TestFirstEmptyLine(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("\r\n+\r\n")

	response, err := call(w, r, "TEST")
	c.Assert(err, IsNil)
	c.Assert(response, IsNil)
}

func (s *CallTestSuite) TestResponseOk(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("+\r\n")

	response, err := call(w, r, "TEST")
	c.Assert(err, IsNil)
	c.Assert(response, IsNil)
}

func (s *CallTestSuite) TestResponseError(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("-MESSAGE\r\n")

	response, err := call(w, r, "TEST")
	c.Assert(err, ErrorMatches, "MESSAGE")
	c.Assert(response, IsNil)
}

func (s *CallTestSuite) TestResponseSingleLine(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("\"value\"\r\n+\r\n")

	response, err := call(w, r, "TEST")
	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, []string{`"value"`})
}

func (s *CallTestSuite) TestResponseMultiLine(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("key1\r\nkey2\r\nkey3\r\n+\r\n")

	response, err := call(w, r, "TEST")
	c.Assert(err, IsNil)
	c.Assert(response, DeepEquals, []string{"key1", "key2", "key3"})
}
