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
		protocol.NewGetRequest().Command(): func(data io.Reader) []byte {
			request := protocol.NewGetRequest()
			response := protocol.NewGetResponse()
			return run(data, request, response, func() {
				response.Value = "value"
			})
		},
		protocol.NewSetRequest().Command(): func(data io.Reader) []byte {
			request := protocol.NewSetRequest()
			response := protocol.NewSetResponse()
			return run(data, request, response, func() {})
		},
		protocol.NewHashGetAllRequest().Command(): func(data io.Reader) []byte {
			request := protocol.NewHashGetAllRequest()
			response := protocol.NewHashGetAllResponse()
			return run(data, request, response, func() {
				response.Fields = map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				}
			})
		},
		protocol.NewListRangeRequest().Command(): func(data io.Reader) []byte {
			request := protocol.NewListRangeRequest()
			response := protocol.NewListRangeResponse()
			return run(data, request, response, func() {
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
	data, _ := r.Encode()
	s.benchmarkCommand(c, data)
}

func (s *SessionTestSuite) BenchmarkSet(c *C) {
	r := protocol.NewSetRequest()
	r.Key = "key"
	r.Value = "value"
	r.TTL = 60
	data, _ := r.Encode()
	s.benchmarkCommand(c, data)
}

func (s *SessionTestSuite) BenchmarkHashGetAll(c *C) {
	r := protocol.NewHashGetAllRequest()
	r.Key = "hash"
	data, _ := r.Encode()
	s.benchmarkCommand(c, data)
}

func (s *SessionTestSuite) BenchmarkListRange(c *C) {
	r := protocol.NewListRangeRequest()
	r.Key = "list"
	data, _ := r.Encode()
	s.benchmarkCommand(c, data)
}

func (s *SessionTestSuite) TestCommand(c *C) {
	commands := map[string]command{
		protocol.NewGetRequest().Command(): func(data io.Reader) []byte {
			request := protocol.NewGetRequest()
			response := protocol.NewGetResponse()
			return run(data, request, response, func() {
				c.Assert(request.Key, Equals, "key")
				response.Value = "value"
			})
		},
	}

	conn := newTestConn()

	go newSession("test", conn, commands, nil, log.New(&bytes.Buffer{}, "", 0)).start()

	request := protocol.NewGetRequest()
	request.Key = "key"
	data, _ := request.Encode()

	out := conn.send(data)
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
