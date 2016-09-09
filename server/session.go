package server

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/Barberrrry/jcache/protocol"
	"github.com/Barberrrry/jcache/server/htpasswd"
)

type session struct {
	id              string
	rwc             io.ReadWriteCloser
	serverCommands  map[string]command
	sessionCommands map[string]command
	isAuthRequired  bool
	isAuthorized    bool
	logger          *log.Logger
}

var (
	unknownCommandError = errors.New("Unknown command")
	needAuthError       = errors.New("Need authentitication")
)

func newSession(id string, rwc io.ReadWriteCloser, commands map[string]command, htpasswdFile *htpasswd.HtpasswdFile, logger *log.Logger) *session {
	s := &session{
		id:             id,
		rwc:            rwc,
		serverCommands: commands,
		logger:         logger,
	}

	if htpasswdFile != nil {
		s.isAuthRequired = true
	}
	s.sessionCommands = map[string]command{
		protocol.NewAuthRequest().Command(): newAuthCommand(htpasswdFile, s),
	}

	return s
}

func (s *session) start() {
	defer s.rwc.Close()

	s.log("open session")
	defer s.log("close session")

	for {
		commandName, err := protocol.ReadRequestCommand(s.rwc)
		if err != nil {
			s.log(fmt.Sprintf("read error: %s", err))
			return
		}

		s.log(fmt.Sprintf("command: %s", commandName))

		commandError := unknownCommandError
		if command, found := s.sessionCommands[commandName]; found {
			command(s.rwc)
			continue
		}

		if command, found := s.serverCommands[commandName]; found {
			if !s.isAuthRequired || s.isAuthorized {
				command(s.rwc)
				continue
			}
			commandError = needAuthError
		}

		s.log(fmt.Sprintf("command error: %s", commandError))
		protocol.FlushRequest(s.rwc)
		writeError(s.rwc, commandError)
	}
}

func (s *session) authorize() {
	s.isAuthorized = true
	s.log("successful authentication")
}

func (s *session) log(message string) {
	s.logger.Printf("[%s] %s", s.id, message)
}
