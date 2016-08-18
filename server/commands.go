package server

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Barberrrry/jcache/server/htpasswd"
)

const (
	keyPattern   = `([a-zA-Z0-9_]+)`
	valuePattern = `"(.*)"`
	fieldPattern = `([a-zA-Z0-9_]+)`
	ttlPattern   = `([a-zA-Z0-9-.]+)`
	intPattern   = `([0-9]+)`

	okResponse             = "OK"
	unknownCommandResponse = "UNKNOWN COMMAND"
	invalidFormatResponse  = "INVALID COMMAND FORMAT"
	needAuthResponse       = "NEED AUTHENTICATION"
	errorResponse          = "ERROR: %s"
	valueResponse          = `"%s"`
	hashElementResponse    = `%s: "%s"`
)

type command struct {
	format     *regexp.Regexp
	run        func(session *session, params []string) string
	allowGuest bool
}

func newKeysCommand() *command {
	return &command{
		format: regexp.MustCompile("^$"),
		run: func(session *session, params []string) string {
			keys := session.server.storage.Keys()
			return strings.Join(keys, "\n")
		},
	}
}

func newTTLCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			ttl, err := session.server.storage.TTL(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return ttl.String()
		},
	}
}

func newGetCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.Get(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newSetCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, valuePattern, ttlPattern)),
		run: func(session *session, params []string) string {
			ttl, err := time.ParseDuration(params[2])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			if err := session.server.storage.Set(params[0], params[1], ttl); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newUpdCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.Update(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newDelCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.Delete(params[0]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashCreateCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(session *session, params []string) string {
			ttl, err := time.ParseDuration(params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			err = session.server.storage.HashCreate(params[0], ttl)
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashGetAllCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			hash, err := session.server.storage.HashGetAll(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			var response []string
			for key, value := range hash {
				response = append(response, fmt.Sprintf(hashElementResponse, key, value))
			}

			return strings.Join(response, "\n")
		},
	}
}

func newHashGetCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.HashGet(params[0], params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newHashSetCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, fieldPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.HashSet(params[0], params[1], params[2]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashDelCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.HashDelete(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashLenCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			len, err := session.server.storage.HashLen(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf("%d", len)
		},
	}
}

func newHashKeysCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			keys, err := session.server.storage.HashKeys(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return strings.Join(keys, "\n")
		},
	}
}

func newListCreateCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(session *session, params []string) string {
			ttl, err := time.ParseDuration(params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			err = session.server.storage.ListCreate(params[0], ttl)
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newListLeftPopCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.ListLeftPop(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newListRightPopCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.ListRightPop(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newListLeftPushCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.ListLeftPush(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newListRightPushCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.ListRightPush(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newListLenCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^LLEN %s$", keyPattern)),
		run: func(session *session, params []string) string {
			len, err := session.server.storage.ListLen(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf("%d", len)
		},
	}
}

func newListRangeCommand() *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, intPattern, intPattern)),
		run: func(session *session, params []string) string {
			start, err := strconv.Atoi(params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			stop, err := strconv.Atoi(params[2])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}

			values, err := session.server.storage.ListRange(params[0], start, stop)
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			var response []string
			for _, value := range values {
				response = append(response, fmt.Sprintf(valueResponse, value))
			}

			return strings.Join(response, "\n")
		},
	}
}

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile) *command {
	return &command{
		allowGuest: true,
		format:     regexp.MustCompile("^([a-zA-Z0-9]+) (.+)$"),
		run: func(session *session, params []string) string {
			if htpasswdFile.Validate(params[0], params[1]) {
				session.authorize()
				return okResponse
			}
			return "INVALID AUTH"
		},
	}
}
