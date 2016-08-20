package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Barberrrry/jcache/server/htpasswd"
)

const (
	lineSeparator = "\r\n"
)

type session struct {
	id             string
	rw             io.ReadWriter
	commands       map[string]*command
	isAuthRequired bool
	isAuthorized   bool
}

func newSession(id string, rw io.ReadWriter, storage storage, htpasswdFile *htpasswd.HtpasswdFile) *session {
	commands := map[string]*command{
		"KEYS":    newKeysCommand(storage),
		"TTL":     newTTLCommand(storage),
		"GET":     newGetCommand(storage),
		"SET":     newSetCommand(storage),
		"DEL":     newDelCommand(storage),
		"UPD":     newUpdCommand(storage),
		"HCREATE": newHashCreateCommand(storage),
		"HGETALL": newHashGetAllCommand(storage),
		"HGET":    newHashGetCommand(storage),
		"HSET":    newHashSetCommand(storage),
		"HDEL":    newHashDelCommand(storage),
		"HLEN":    newHashLenCommand(storage),
		"HKEYS":   newHashKeysCommand(storage),
		"LCREATE": newListCreateCommand(storage),
		"LLPOP":   newListLeftPopCommand(storage),
		"LRPOP":   newListRightPopCommand(storage),
		"LLPUSH":  newListLeftPushCommand(storage),
		"LRPUSH":  newListRightPushCommand(storage),
		"LLEN":    newListLenCommand(storage),
		"LRANGE":  newListRangeCommand(storage),
	}

	s := &session{
		id:       id,
		rw:       rw,
		commands: commands,
	}

	if htpasswdFile != nil {
		s.isAuthRequired = true
	}
	s.commands["AUTH"] = newAuthCommand(htpasswdFile, s)

	return s
}

func (s *session) start() {
	s.log("start session")
	for {
		if err := s.handleInput(); err != nil {
			break
		}
	}
	s.log("close session")
}

func (s *session) handleInput() error {
	buf := bufio.NewReader(s.rw)
	line, _, err := buf.ReadLine()

	if err != nil {
		return err
	}

	if len(line) > 0 {
		response := s.handleCommand(string(line))
		for _, l := range response {
			s.rw.Write([]byte(l + lineSeparator))
		}
	}

	return nil
}

func (s *session) handleCommand(line string) []string {
	parts := strings.SplitN(line, " ", 2)
	if command, found := s.commands[parts[0]]; found {
		if s.isAuthRequired && !command.allowGuest && !s.isAuthorized {
			return errorResponse(errors.New("NEED AUTHENTICATION"))
		}

		var arguments string
		if len(parts) > 1 {
			arguments = parts[1]
		}
		matches := command.format.FindStringSubmatch(arguments)
		if len(matches) > 0 {
			var params []string
			if len(matches) > 1 {
				params = matches[1:]
			}
			s.log(fmt.Sprintf("run %s", parts[0]))
			return command.run(params)
		}
		return errorResponse(errors.New("INVALID COMMAND FORMAT"))
	}
	return errorResponse(errors.New("UNKNOWN COMMAND"))
}

func (s *session) authorize() {
	s.isAuthorized = true
	s.log("successful authentication")
}

func (s *session) log(message string) {
	log.Printf("[%s] %s", s.id, message)
}
