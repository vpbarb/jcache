package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/storage"
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

func newSession(id string, rw io.ReadWriter, storage storage.Storage, htpasswdFile *htpasswd.HtpasswdFile) *session {
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
	rb := bufio.NewReader(s.rw)
	for {
		line, _, err := rb.ReadLine()
		if err != nil {
			s.log(fmt.Sprintf("read error: %s", err))
			break
		}

		if len(line) > 0 {
			parts := strings.SplitN(string(line), " ", 2)
			if command, found := s.commands[parts[0]]; found {
				if s.isAuthRequired && !command.allowGuest && !s.isAuthorized {
					s.writeResponse(errorResponse("NEED AUTHENTICATION"))
					continue
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
					s.writeResponse(command.run(params, rb))
				} else {
					s.writeResponse(errorResponse("INVALID COMMAND FORMAT"))
				}
			} else {
				s.writeResponse(errorResponse("UNKNOWN COMMAND"))
			}
		}
	}
	s.log("close session")
}

func (s *session) writeResponse(response []string) {
	for _, line := range response {
		s.rw.Write([]byte(line + lineSeparator))
	}
}

func (s *session) authorize() {
	s.isAuthorized = true
	s.log("successful authentication")
}

func (s *session) log(message string) {
	log.Printf("[%s] %s", s.id, message)
}
