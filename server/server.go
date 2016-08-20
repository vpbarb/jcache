package server

import (
	"log"
	"net"

	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/storage"
)

type server struct {
	storage      storage.Storage
	htpasswdFile *htpasswd.HtpasswdFile
}

func New(storage storage.Storage, htpasswdPath string) *server {
	s := &server{storage: storage}

	if htpasswdPath != "" {
		var err error
		if s.htpasswdFile, err = htpasswd.NewHtpasswdFromFile(htpasswdPath); err == nil {
			log.Print("server supports authentication")
		} else {
			log.Printf("erron on loading htpasswd file: %s", err)
		}
	}

	return s
}

func (s *server) ListenAndServe(addr string) {
	log.Printf("listen on %s", addr)
	listener, _ := net.Listen("tcp", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error on connection accepting: %s\n", err)
			continue
		}

		go newSession(conn.RemoteAddr().String(), conn, s.storage, s.htpasswdFile).start()
	}
}
