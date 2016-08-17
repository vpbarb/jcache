package server

import (
	"log"

	"github.com/Barberrrry/jcache/server/transport"
)

type Server struct {
}

func (s *Server) Handle(c *transport.Request) transport.Response {
	log.Printf("%+v", c)

	if c.Name == "error" {
		return []byte("TEST ERROR")
	}

	return nil
}
