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

type session struct {
	id             string
	rw             io.ReadWriter
	commands       map[string]command
	authCommand    command
	isAuthRequired bool
	isAuthorized   bool
}

func newSession(id string, rw io.ReadWriter, commands map[string]command, htpasswdFile *htpasswd.HtpasswdFile) *session {
	s := &session{
		id:       id,
		rw:       rw,
		commands: commands,
	}

	if htpasswdFile != nil {
		s.isAuthRequired = true
	}
	s.authCommand = newAuthCommand(htpasswdFile, s)

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

			if parts[0] == "AUTH" {
				s.log(fmt.Sprintf("run %s", parts[0]))
				s.rw.Write(s.authCommand(line, rb))
				continue
			}

			if command, found := s.commands[parts[0]]; found {
				if s.isAuthRequired && !s.isAuthorized {
					s.rw.Write(formatError(errors.New("Need authentitication")))
					continue
				}
				s.log(fmt.Sprintf("run %s", parts[0]))
				s.rw.Write(command(line, rb))
			} else {
				s.rw.Write(formatError(errors.New("Unknown command")))
			}
		}
	}
	s.log("close session")
}

func (s *session) authorize() {
	s.isAuthorized = true
	s.log("successful authentication")
}

func (s *session) log(message string) {
	log.Printf("[%s] %s", s.id, message)
}
