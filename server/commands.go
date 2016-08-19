package server

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/Barberrrry/jcache/server/htpasswd"
)

const (
	keyPattern   = `([a-zA-Z0-9_]+)`
	valuePattern = `"(.*)"`
	fieldPattern = `([a-zA-Z0-9_]+)`
	ttlPattern   = `([a-zA-Z0-9-.]+)`
	intPattern   = `([0-9]+)`

	errorTemplate       = "ERROR: %s\r\n"
	valueTemplate       = "\"%s\"\r\n"
	keyTemplate         = "%s\r\n"
	countTemplate       = "(%d)\r\n"
	hashElementTemplate = "%s: \"%s\"\r\n"
)

var (
	okResponse             = []byte("OK\r\n")
	unknownCommandResponse = []byte("UNKNOWN COMMAND\r\n")
	invalidFormatResponse  = []byte("INVALID COMMAND FORMAT\r\n")
	needAuthResponse       = []byte("NEED AUTHENTICATION\r\n")
)

type command struct {
	format     *regexp.Regexp
	run        func(params []string) []byte
	allowGuest bool
}

func errorResponse(err error) []byte {
	return []byte(fmt.Sprintf(errorTemplate, err))
}

func valueResponse(value string) []byte {
	return []byte(fmt.Sprintf(valueTemplate, value))
}

func newKeysCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile("^$"),
		run: func(params []string) []byte {
			keys := storage.Keys()
			response := bytes.NewBufferString(fmt.Sprintf(countTemplate, len(keys)))
			for _, key := range keys {
				response.WriteString(fmt.Sprintf(keyTemplate, key))
			}
			return response.Bytes()
		},
	}
}

func newTTLCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			ttl, err := storage.TTL(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return []byte(ttl.String())
		},
	}
}

func newGetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			value, err := storage.Get(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newSetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, valuePattern, ttlPattern)),
		run: func(params []string) []byte {
			ttl, err := time.ParseDuration(params[2])
			if err != nil {
				return errorResponse(err)
			}
			if err := storage.Set(params[0], params[1], ttl); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newUpdCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) []byte {
			if err := storage.Update(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newDelCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			if err := storage.Delete(params[0]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashCreateCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(params []string) []byte {
			ttl, err := time.ParseDuration(params[1])
			if err != nil {
				return errorResponse(err)
			}
			err = storage.HashCreate(params[0], ttl)
			if err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashGetAllCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			hash, err := storage.HashGetAll(params[0])
			if err != nil {
				return errorResponse(err)
			}
			response := bytes.NewBufferString(fmt.Sprintf(countTemplate, len(hash)))
			for key, value := range hash {
				response.WriteString(fmt.Sprintf(hashElementTemplate, key, value))
			}

			return response.Bytes()
		},
	}
}

func newHashGetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(params []string) []byte {
			value, err := storage.HashGet(params[0], params[1])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newHashSetCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, fieldPattern, valuePattern)),
		run: func(params []string) []byte {
			if err := storage.HashSet(params[0], params[1], params[2]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashDelCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(params []string) []byte {
			if err := storage.HashDelete(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashLenCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			len, err := storage.HashLen(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return []byte(fmt.Sprintf("%d", len))
		},
	}
}

func newHashKeysCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			keys, err := storage.HashKeys(params[0])
			if err != nil {
				return errorResponse(err)
			}
			response := bytes.NewBufferString(fmt.Sprintf(countTemplate, len(keys)))
			for _, key := range keys {
				response.WriteString(fmt.Sprintf(keyTemplate, key))
			}
			return response.Bytes()
		},
	}
}

func newListCreateCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(params []string) []byte {
			ttl, err := time.ParseDuration(params[1])
			if err != nil {
				return errorResponse(err)
			}
			err = storage.ListCreate(params[0], ttl)
			if err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newListLeftPopCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			value, err := storage.ListLeftPop(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newListRightPopCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			value, err := storage.ListRightPop(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newListLeftPushCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) []byte {
			if err := storage.ListLeftPush(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newListRightPushCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) []byte {
			if err := storage.ListRightPush(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newListLenCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []byte {
			len, err := storage.ListLen(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return []byte(fmt.Sprintf("%d", len))
		},
	}
}

func newListRangeCommand(storage storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, intPattern, intPattern)),
		run: func(params []string) []byte {
			start, err := strconv.Atoi(params[1])
			if err != nil {
				return errorResponse(err)
			}
			stop, err := strconv.Atoi(params[2])
			if err != nil {
				return errorResponse(err)
			}

			values, err := storage.ListRange(params[0], start, stop)
			if err != nil {
				return errorResponse(err)
			}
			response := bytes.NewBufferString(fmt.Sprintf(countTemplate, len(values)))
			for _, value := range values {
				response.Write(valueResponse(value))
			}

			return response.Bytes()
		},
	}
}

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile, session *session) *command {
	return &command{
		allowGuest: true,
		format:     regexp.MustCompile("^([a-zA-Z0-9]+) (.+)$"),
		run: func(params []string) []byte {
			if htpasswdFile == nil || htpasswdFile.Validate(params[0], params[1]) {
				session.authorize()
				return okResponse
			}
			return []byte("INVALID AUTH\r\n")
		},
	}
}
