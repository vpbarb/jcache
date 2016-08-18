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
	run        func(params []string) string
	allowGuest bool
}

func newKeysCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile("^$"),
		run: func(params []string) string {
			keys := storage.Keys()
			return strings.Join(keys, "\n")
		},
	}
}

func newTTLCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			ttl, err := storage.TTL(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return ttl.String()
		},
	}
}

func newGetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			value, err := storage.Get(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newSetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, valuePattern, ttlPattern)),
		run: func(params []string) string {
			ttl, err := time.ParseDuration(params[2])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			if err := storage.Set(params[0], params[1], ttl); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newUpdCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) string {
			if err := storage.Update(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newDelCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			if err := storage.Delete(params[0]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashCreateCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(params []string) string {
			ttl, err := time.ParseDuration(params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			err = storage.HashCreate(params[0], ttl)
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashGetAllCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			hash, err := storage.HashGetAll(params[0])
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

func newHashGetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(params []string) string {
			value, err := storage.HashGet(params[0], params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newHashSetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, fieldPattern, valuePattern)),
		run: func(params []string) string {
			if err := storage.HashSet(params[0], params[1], params[2]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashDelCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(params []string) string {
			if err := storage.HashDelete(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newHashLenCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			len, err := storage.HashLen(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf("%d", len)
		},
	}
}

func newHashKeysCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			keys, err := storage.HashKeys(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return strings.Join(keys, "\n")
		},
	}
}

func newListCreateCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(params []string) string {
			ttl, err := time.ParseDuration(params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			err = storage.ListCreate(params[0], ttl)
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newListLeftPopCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			value, err := storage.ListLeftPop(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newListRightPopCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) string {
			value, err := storage.ListRightPop(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}
}

func newListLeftPushCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) string {
			if err := storage.ListLeftPush(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newListRightPushCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) string {
			if err := storage.ListRightPush(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}
}

func newListLenCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^LLEN %s$", keyPattern)),
		run: func(params []string) string {
			len, err := storage.ListLen(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf("%d", len)
		},
	}
}

func newListRangeCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, intPattern, intPattern)),
		run: func(params []string) string {
			start, err := strconv.Atoi(params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			stop, err := strconv.Atoi(params[2])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}

			values, err := storage.ListRange(params[0], start, stop)
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

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile, session *session) *command {
	return &command{
		allowGuest: true,
		format:     regexp.MustCompile("^([a-zA-Z0-9]+) (.+)$"),
		run: func(params []string) string {
			if htpasswdFile.Validate(params[0], params[1]) {
				session.authorize()
				return okResponse
			}
			return "INVALID AUTH"
		},
	}
}
