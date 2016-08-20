package server

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/storage"
)

const (
	keyPattern   = `([a-zA-Z0-9_]+)`
	valuePattern = `"(.*)"`
	fieldPattern = `([a-zA-Z0-9_]+)`
	ttlPattern   = `([a-zA-Z0-9-.]+)`
	intPattern   = `([0-9]+)`

	errorTemplate       = "-%s"
	keyTemplate         = "%s"
	valueTemplate       = `"%s"`
	hashElementTemplate = `%s:"%s"`
	lenTemplate         = "%d"
	okTemplate          = "+"
)

var (
	okResponse = []string{okTemplate}
)

type command struct {
	format     *regexp.Regexp
	run        func(params []string) []string
	allowGuest bool
}

func errorResponse(err error) []string {
	return []string{fmt.Sprintf(errorTemplate, err)}
}

func valueResponse(value string) []string {
	return []string{fmt.Sprintf(valueTemplate, value), okTemplate}
}

func lenResponse(len int) []string {
	return []string{fmt.Sprintf(lenTemplate, len), okTemplate}
}

func newKeysCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile("^$"),
		run: func(params []string) []string {
			keys := storage.Keys()
			var response []string
			for _, key := range keys {
				response = append(response, fmt.Sprintf(keyTemplate, key))
			}
			response = append(response, okTemplate)
			return response
		},
	}
}

func newTTLCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			ttl, err := storage.TTL(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return []string{ttl.String(), okTemplate}
		},
	}
}

func newGetCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			value, err := storage.Get(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newSetCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, valuePattern, ttlPattern)),
		run: func(params []string) []string {
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

func newUpdCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) []string {
			if err := storage.Update(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newDelCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			if err := storage.Delete(params[0]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashCreateCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(params []string) []string {
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

func newHashGetAllCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			hash, err := storage.HashGetAll(params[0])
			if err != nil {
				return errorResponse(err)
			}
			var response []string
			for key, value := range hash {
				response = append(response, fmt.Sprintf(hashElementTemplate, key, value))
			}
			response = append(response, okTemplate)
			return response
		},
	}
}

func newHashGetCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(params []string) []string {
			value, err := storage.HashGet(params[0], params[1])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newHashSetCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, fieldPattern, valuePattern)),
		run: func(params []string) []string {
			if err := storage.HashSet(params[0], params[1], params[2]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashDelCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(params []string) []string {
			if err := storage.HashDelete(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashLenCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			len, err := storage.HashLen(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return lenResponse(len)
		},
	}
}

func newHashKeysCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			keys, err := storage.HashKeys(params[0])
			if err != nil {
				return errorResponse(err)
			}
			var response []string
			for _, key := range keys {
				response = append(response, fmt.Sprintf(keyTemplate, key))
			}
			response = append(response, okTemplate)
			return response
		},
	}
}

func newListCreateCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, ttlPattern)),
		run: func(params []string) []string {
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

func newListLeftPopCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			value, err := storage.ListLeftPop(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newListRightPopCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			value, err := storage.ListRightPop(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return valueResponse(value)
		},
	}
}

func newListLeftPushCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) []string {
			if err := storage.ListLeftPush(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newListRightPushCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(params []string) []string {
			if err := storage.ListRightPush(params[0], params[1]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newListLenCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string) []string {
			len, err := storage.ListLen(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return lenResponse(len)
		},
	}
}

func newListRangeCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, intPattern, intPattern)),
		run: func(params []string) []string {
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
			var response []string
			for _, value := range values {
				response = append(response, fmt.Sprintf(valueTemplate, value))
			}
			response = append(response, okTemplate)
			return response
		},
	}
}

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile, session *session) *command {
	return &command{
		allowGuest: true,
		format:     regexp.MustCompile("^([a-zA-Z0-9]+) (.+)$"),
		run: func(params []string) []string {
			if htpasswdFile == nil || htpasswdFile.Validate(params[0], params[1]) {
				session.authorize()
				return okResponse
			}
			return errorResponse(errors.New("INVALID CREDENTIALS"))
		},
	}
}
