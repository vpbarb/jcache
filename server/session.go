package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type session struct {
	conn         net.Conn
	server       *server
	isAuthorized bool
	remoteAddr   string
}

func (s *session) serve() {
	s.log("start session")
	for {
		if err := s.handleInput(); err != nil {
			break
		}
	}
	s.log("close session")
}

func (s *session) handleInput() error {
	buf := bufio.NewReader(s.conn)
	line, _, err := buf.ReadLine()

	if err != nil {
		return err
	}

	if len(line) > 0 {
		if response := s.handleCommand(string(line)); len(response) > 0 {
			s.conn.Write([]byte(response))
			s.conn.Write([]byte("\n"))
		}
	}

	return nil
}

func (s *session) handleCommand(line string) string {
	parts := strings.SplitN(line, " ", 2)
	if command, found := s.server.commands[parts[0]]; found {
		if s.server.isAuthEnabled && !command.allowGuest && !s.isAuthorized {
			return needAuthResponse
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
			return command.run(s, params)
		}
		return invalidFormatResponse
	}
	return unknownCommandResponse
}

func (s *session) authorize() {
	s.isAuthorized = true
	s.log("successful authentication")
}

func (s *session) log(message string) {
	log.Printf("[%s] %s", s.remoteAddr, message)
}
