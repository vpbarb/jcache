package server

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/Barberrrry/jcache/protocol"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type SessionTestSuite struct{}

var _ = Suite(&SessionTestSuite{})

func (s *SessionTestSuite) TestCommand(c *C) {
	commands := map[string]command{
		protocol.NewGetRequest().Command(): func(rw io.ReadWriter) {
			request := protocol.NewGetRequest()
			response := protocol.NewGetResponse()
			run(rw, request, response, func() {
				c.Assert(request.Key, Equals, "key")
				response.Value = "value"
			})
		},
	}

	conn := newTestConn()

	go newSession("test", conn, commands, nil, log.New(&bytes.Buffer{}, "", 0)).start()

	request := protocol.NewGetRequest()
	request.Key = "key"
	request.Encode(conn.inWriter)

	response := protocol.NewGetResponse()
	err := response.Decode(conn.outReader)
	c.Assert(err, IsNil)
	c.Assert(response.Value, Equals, "value")

	conn.inWriter.Close()
}

type testConn struct {
	inReader  *io.PipeReader
	inWriter  *io.PipeWriter
	outReader *io.PipeReader
	outWriter *io.PipeWriter
}

func newTestConn() *testConn {
	c := &testConn{}
	c.inReader, c.inWriter = io.Pipe()
	c.outReader, c.outWriter = io.Pipe()
	return c
}

func (c *testConn) Read(p []byte) (int, error) {
	return c.inReader.Read(p)
}

func (c *testConn) Write(p []byte) (int, error) {
	return c.outWriter.Write(p)
}
