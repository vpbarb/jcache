package protocol

import (
	"testing"

	"bytes"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ProtocolTestSuite struct{}

var _ = Suite(&ProtocolTestSuite{})

func (s *ProtocolTestSuite) BenchmarkGetRequestEncode(c *C) {
	r := NewGetRequest()
	r.Key = "key"

	w := &bytes.Buffer{}

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Encode(w)
	}
}

func (s *ProtocolTestSuite) BenchmarkSetRequestEncode(c *C) {
	r := NewSetRequest()
	r.Key = "key"
	r.Value = "value"
	r.TTL = 60

	w := &bytes.Buffer{}

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Encode(w)
	}
}

func (s *ProtocolTestSuite) BenchmarkGetResponseEncode(c *C) {
	r := NewGetResponse()
	r.Value = "value"

	w := &bytes.Buffer{}

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Encode(w)
	}
}

func (s *ProtocolTestSuite) BenchmarkSetResponseEncode(c *C) {
	r := NewSetResponse()

	w := &bytes.Buffer{}

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Encode(w)
	}
}

func (s *ProtocolTestSuite) BenchmarkGetRequestDecode(c *C) {
	r := NewGetRequest()

	reader := bytes.NewBufferString("GET key\r\n")

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Decode(reader)
	}
}

func (s *ProtocolTestSuite) BenchmarkSetRequestDecode(c *C) {
	r := NewSetRequest()
	r.Key = "key"
	r.Value = "value"
	r.TTL = 60

	reader := bytes.NewBufferString("SET key 60 5\r\nvalue\r\n")

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Decode(reader)
	}
}

func (s *ProtocolTestSuite) BenchmarkGetResponseDecode(c *C) {
	r := NewGetResponse()
	r.Value = "value"

	reader := bytes.NewBufferString("VALUE 5\r\nvalue\r\n")

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Decode(reader)
	}
}

func (s *ProtocolTestSuite) BenchmarkSetResponseDecode(c *C) {
	r := NewSetResponse()

	reader := bytes.NewBufferString("OK\r\n")

	c.ResetTimer()
	for i := 0; i < c.N; i++ {
		r.Decode(reader)
	}
}
