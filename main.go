package main

import (
	"github.com/Barberrrry/jcache/server"
	"github.com/Barberrrry/jcache/server/transport"
)

func main() {
	s := &server.Server{}
	transport.ListenTCP(s, ":8081")
}
