package server

import "fmt"

var (
	okResponse             = "OK"
	unknownCommandResponse = "UNKNOWN COMMAND"
	invalidFormatResponse  = "INVALID FORMAT"
	errorResponse          = "ERROR: %s"
)

type getCommand struct{}

func (c *getCommand) run(session *session, params []string) string {
	if len(params) < 1 {
		return invalidFormatResponse
	}
	value, err := session.server.storage.Get(params[0])
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return value
}

type setCommand struct{}

func (c *setCommand) run(session *session, params []string) string {
	if len(params) < 2 {
		return invalidFormatResponse
	}
	if err := session.server.storage.Set(params[0], params[1]); err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return okResponse
}

type delCommand struct{}

func (c *delCommand) run(session *session, params []string) string {
	if len(params) < 1 {
		return invalidFormatResponse
	}
	if err := session.server.storage.Delete(params[0]); err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return okResponse
}
