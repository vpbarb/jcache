package server

import (
	"bufio"
	"github.com/Barberrrry/jcache/server/memory"
	"log"
	"net"
	"strings"
	"time"
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

type commandFunc func(session *session, params string) string

type server struct {
	storage  storage
	commands map[string]commandFunc
}

type session struct {
	conn         net.Conn
	server       *server
	isAuthorized bool
}

func New() *server {
	return &server{
		storage: memory.NewStorage(),
		commands: map[string]commandFunc{
			"KEYS":    keysCommand,
			"TTL":     ttlCommand,
			"GET":     getCommand,
			"SET":     setCommand,
			"DEL":     delCommand,
			"UPD":     updCommand,
			"HCREATE": hashCreateCommand,
			"HGETALL": hashGetAllCommand,
			"HGET":    hashGetCommand,
			"HSET":    hashSetCommand,
			"HDEL":    hashDelCommand,
			"HLEN":    hashLenCommand,
			"HKEYS":   hashKeysCommand,
			//"LCREATE": listCreateCommand,
			//"LLPOP":   listLeftPopCommand,
			//"LRPOP":   listRightPopCommand,
			//"LLPUSH":  listLeftPushCommand,
			//"LRPUSH":  listRightPushCommand,
			//"LLEN":    listLenCommand,
			//"LRANGE":  listRangeCommand,
		},
	}
}

func (s *server) ListenAndServe(addr string) {
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
		conn:   conn,
		server: s,
	}
}

func (s *session) serve() {
	log.Print("start session")
	for {
		if err := s.handleCommand(); err != nil {
			break
		}
	}
	log.Print("close session")
}

func (s *session) handleCommand() error {
	buf := bufio.NewReader(s.conn)
	line, _, err := buf.ReadLine()

	if err != nil {
		return err
	}

	if len(line) > 0 {
		parts := strings.SplitN(string(line), " ", 2)

		var response string
		if command, found := s.server.commands[parts[0]]; found {
			var params string
			if len(parts) > 1 {
				params = parts[1]
			}
			response = command(s, params)
		} else {
			response = unknownCommandResponse
		}

		if len(response) > 0 {
			s.conn.Write([]byte(response))
			s.conn.Write([]byte("\n"))
		}
	}

	return nil
}
