package client

import (
	"bufio"
	"bytes"
	"errors"

	. "gopkg.in/check.v1"
)

type CallTestSuite struct{}

var (
	_ = Suite(&CallTestSuite{})

	nilDataFormatter = func(response *bufio.Reader) (err error) {
		return nil
	}
)

type errorWriter struct {
	err error
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, w.err
}

func (s *CallTestSuite) TestWriteError(c *C) {
	w := &errorWriter{errors.New("Write error")}
	r := &bytes.Buffer{}

	err := transfer(w, r, "TEST\r\n", nilDataFormatter)
	c.Assert(err, ErrorMatches, "Cannot write to connection: Write error")
}

func (s *CallTestSuite) TestWriteOk(c *C) {
	w := &bytes.Buffer{}
	r := &bytes.Buffer{}

	transfer(w, r, "TEST\r\n", nilDataFormatter)
	c.Assert(w.String(), DeepEquals, "TEST\r\n")
}

func (s *CallTestSuite) TestReadError(c *C) {
	w := &bytes.Buffer{}
	r := &bytes.Buffer{}

	err := transfer(w, r, "TEST\r\n", nilDataFormatter)
	c.Assert(err, ErrorMatches, "Cannot read from connection: EOF")
}

func (s *CallTestSuite) TestResponseOk(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("OK\r\n")

	err := transfer(w, r, "TEST\r\n", nilDataFormatter)
	c.Assert(err, IsNil)
}

func (s *CallTestSuite) TestResponseError(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("ERROR MESSAGE\r\n")

	err := transfer(w, r, "TEST\r\n", nilDataFormatter)
	c.Assert(err, ErrorMatches, "Response error: MESSAGE")
}

func (s *CallTestSuite) TestResponseData(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("DATA\r\nsome data\r\nEND\r\n")

	var data string
	dataFormatter := func(response *bufio.Reader) (err error) {
		for {
			line, _, err := response.ReadLine()
			if err != nil {
				return err
			}
			str := string(line)
			if str == endResponse {
				return nil
			}
			data = str
		}
		return nil
	}

	err := transfer(w, r, "TEST", dataFormatter)
	c.Assert(err, IsNil)
	c.Assert(data, Equals, "some data")
}

func (s *CallTestSuite) TestResponseInvalid(c *C) {
	w := &bytes.Buffer{}
	r := bytes.NewBufferString("WTF\r\n")

	err := transfer(w, r, "TEST\r\n", nilDataFormatter)
	c.Assert(err, ErrorMatches, "Invalid response format")
}
