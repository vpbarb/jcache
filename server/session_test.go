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

func (s *SessionTestSuite) prepareDummyCommands() map[string]command {
	return map[string]command{
		protocol.NewGetRequest().Command(): func(rw io.ReadWriter) {
			request := protocol.NewGetRequest()
			response := protocol.NewGetResponse()
			run(rw, request, response, func() {
				response.Value = "value"
			})
		},
		protocol.NewSetRequest().Command(): func(rw io.ReadWriter) {
			request := protocol.NewSetRequest()
			response := protocol.NewSetResponse()
			run(rw, request, response, func() {})
		},
		protocol.NewHashGetAllRequest().Command(): func(rw io.ReadWriter) {
			request := protocol.NewHashGetAllRequest()
			response := protocol.NewHashGetAllResponse()
			run(rw, request, response, func() {
				response.Fields = map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				}
			})
		},
		protocol.NewListRangeRequest().Command(): func(rw io.ReadWriter) {
			request := protocol.NewListRangeRequest()
			response := protocol.NewListRangeResponse()
			run(rw, request, response, func() {
				response.Values = []string{
					"value1",
					"value2",
					"value3",
				}
			})
		},
	}
}

func (s *SessionTestSuite) benchmarkCommand(c *C, data []byte) {
	commands := s.prepareDummyCommands()

	conn := newTestConn()

	go newSession("test", conn, commands, nil, log.New(&bytes.Buffer{}, "", 0)).start()

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		conn.send(data)
	}
	conn.close()
}

func (s *SessionTestSuite) BenchmarkGet(c *C) {
	r := protocol.NewGetRequest()
	r.Key = "key"
	data := &bytes.Buffer{}
	r.Encode(data)
	s.benchmarkCommand(c, data.Bytes())
}

func (s *SessionTestSuite) BenchmarkSet(c *C) {
	r := protocol.NewSetRequest()
	r.Key = "key"
	r.Value = "value"
	r.TTL = 60
	data := &bytes.Buffer{}
	r.Encode(data)
	s.benchmarkCommand(c, data.Bytes())
}

func (s *SessionTestSuite) BenchmarkHashGetAll(c *C) {
	r := protocol.NewHashGetAllRequest()
	r.Key = "hash"
	data := &bytes.Buffer{}
	r.Encode(data)
	s.benchmarkCommand(c, data.Bytes())
}

func (s *SessionTestSuite) BenchmarkListRange(c *C) {
	r := protocol.NewListRangeRequest()
	r.Key = "list"
	data := &bytes.Buffer{}
	r.Encode(data)
	s.benchmarkCommand(c, data.Bytes())
}

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
	data := &bytes.Buffer{}
	request.Encode(data)
	out := conn.send(data.Bytes())
	conn.close()

	response := protocol.NewGetResponse()
	err := response.Decode(bytes.NewBuffer(out))
	c.Assert(err, IsNil)
	c.Assert(response.Value, Equals, "value")
}

type testConn struct {
	in  chan byte
	out []byte
}

func newTestConn() *testConn {
	return &testConn{
		in: make(chan byte, 0),
	}
}

func (c *testConn) send(p []byte) []byte {
	for _, b := range p {
		c.in <- b
	}
	return c.out
}

func (c *testConn) close() {
	close(c.in)
}

func (c *testConn) Read(p []byte) (int, error) {
	i := 0
	for b := range c.in {
		p[i] = b
		i++
		if i == len(p) {
			return i, nil
		}
	}
	return i, io.EOF
}

func (c *testConn) Write(p []byte) (int, error) {
	c.out = p
	return len(p), nil
}
