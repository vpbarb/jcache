package server

import (
	"log"
	"net"
	"time"

	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/memory"
)

type storage interface {
	Keys() []string
	TTL(key string) (time.Duration, error)
	Get(key string) (string, error)
	Set(key, value string, ttl time.Duration) error
	Update(key, value string) error
	Delete(key string) error
	HashCreate(key string, ttl time.Duration) error
	HashGet(key, field string) (string, error)
	HashGetAll(key string) (map[string]string, error)
	HashSet(key, field, value string) error
	HashDelete(key, field string) error
	HashLen(key string) (int, error)
	HashKeys(key string) ([]string, error)
	ListCreate(key string, ttl time.Duration) error
	ListLeftPop(key string) (string, error)
	ListRightPop(key string) (string, error)
	ListLeftPush(key, value string) error
	ListRightPush(key, value string) error
	ListLen(key string) (int, error)
	ListRange(key string, start, stop int) ([]string, error)
}

type server struct {
	storage       storage
	commands      map[string]*command
	isAuthEnabled bool
}

func New(htpasswdPath string) *server {
	s := &server{
		storage: memory.NewStorage(),
		commands: map[string]*command{
			"KEYS":    newKeysCommand(),
			"TTL":     newTTLCommand(),
			"GET":     newGetCommand(),
			"SET":     newSetCommand(),
			"DEL":     newDelCommand(),
			"UPD":     newUpdCommand(),
			"HCREATE": newHashCreateCommand(),
			"HGETALL": newHashGetAllCommand(),
			"HGET":    newHashGetCommand(),
			"HSET":    newHashSetCommand(),
			"HDEL":    newHashDelCommand(),
			"HLEN":    newHashLenCommand(),
			"HKEYS":   newHashKeysCommand(),
			"LCREATE": newListCreateCommand(),
			"LLPOP":   newListLeftPopCommand(),
			"LRPOP":   newListRightPopCommand(),
			"LLPUSH":  newListLeftPushCommand(),
			"LRPUSH":  newListRightPushCommand(),
			"LLEN":    newListLenCommand(),
			"LRANGE":  newListRangeCommand(),
		},
	}

	if htpasswdPath != "" {
		if htpasswdFile, err := htpasswd.NewHtpasswdFromFile(htpasswdPath); err == nil {
			s.isAuthEnabled = true
			s.commands["AUTH"] = newAuthCommand(htpasswdFile)
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

		go s.newSession(conn).serve()
	}
}

func (s *server) newSession(conn net.Conn) *session {
	return &session{
		conn:       conn,
		server:     s,
		remoteAddr: conn.RemoteAddr().String(),
	}
}
