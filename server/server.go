package server

import (
	"log"
	"net"

	"github.com/Barberrrry/jcache/protocol"
	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/storage"
)

type server struct {
	commands     map[string]command
	storage      storage.Storage
	htpasswdFile *htpasswd.HtpasswdFile
	logger       *log.Logger
}

func New(storage storage.Storage, htpasswdPath string, logger *log.Logger) *server {
	s := &server{
		storage: storage,
		commands: map[string]command{
			protocol.NewKeysRequest().Command():          newKeysCommand(storage),
			protocol.NewGetRequest().Command():           newGetCommand(storage),
			protocol.NewSetRequest().Command():           newSetCommand(storage),
			protocol.NewDelRequest().Command():           newDelCommand(storage),
			protocol.NewUpdRequest().Command():           newUpdCommand(storage),
			protocol.NewHashCreateRequest().Command():    newHashCreateCommand(storage),
			protocol.NewHashGetAllRequest().Command():    newHashGetAllCommand(storage),
			protocol.NewHashGetRequest().Command():       newHashGetCommand(storage),
			protocol.NewHashSetRequest().Command():       newHashSetCommand(storage),
			protocol.NewHashDelRequest().Command():       newHashDelCommand(storage),
			protocol.NewHashLenRequest().Command():       newHashLenCommand(storage),
			protocol.NewHashKeysRequest().Command():      newHashKeysCommand(storage),
			protocol.NewListCreateRequest().Command():    newListCreateCommand(storage),
			protocol.NewListLeftPopRequest().Command():   newListLeftPopCommand(storage),
			protocol.NewListRightPopRequest().Command():  newListRightPopCommand(storage),
			protocol.NewListLeftPushRequest().Command():  newListLeftPushCommand(storage),
			protocol.NewListRightPushRequest().Command(): newListRightPushCommand(storage),
			protocol.NewListLenRequest().Command():       newListLenCommand(storage),
			protocol.NewListRangeRequest().Command():     newListRangeCommand(storage),
			protocol.NewExpireRequest().Command():        newExpireCommand(storage),
		},
		logger: logger,
	}

	if htpasswdPath != "" {
		var err error
		if s.htpasswdFile, err = htpasswd.NewHtpasswdFromFile(htpasswdPath); err == nil {
			s.logger.Print("server supports authentication")
		} else {
			s.logger.Printf("erron on loading htpasswd file: %s", err)
		}
	}

	return s
}

func (s *server) ListenAndServe(addr string) {
	s.logger.Printf("listen on %s", addr)
	listener, _ := net.Listen("tcp", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error on connection accepting: %s\n", err)
			continue
		}

		go newSession(conn.RemoteAddr().String(), conn, s.commands, s.htpasswdFile, s.logger).start()
	}
}
