package server

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/Barberrrry/jcache/server/htpasswd"
	"github.com/Barberrrry/jcache/server/storage"
)

const (
	keyPattern = `([a-zA-Z0-9_]+)`
	intPattern = `([0-9]+)`

	okTemplate        = "OK"
	errorTemplate     = "ERROR %s"
	dataTemplate      = "DATA"
	endTemplate       = "END"
	keyTemplate       = "KEY %s"
	valueTemplate     = "VALUE %d"
	hashFieldTemplate = "FIELD %s %d"
	lenTemplate       = "LEN %d"
	ttlTemplate       = "TTL %d"
)

var (
	okResponse = []string{okTemplate}
)

type command struct {
	format       *regexp.Regexp
	run          func(params []string, data io.Reader) []string
	requiredData bool
	allowGuest   bool
}

func errorResponse(err interface{}) []string {
	return []string{fmt.Sprintf(errorTemplate, err)}
}

func valueResponse(value string) []string {
	response := []string{dataTemplate}
	response = append(response, valueLines(value)...)
	response = append(response, endTemplate)
	return response
}

func valueLines(value string) []string {
	return []string{fmt.Sprintf(valueTemplate, len(value)), value}
}

func lengthResponse(length int) []string {
	return []string{dataTemplate, fmt.Sprintf(lenTemplate, length), endTemplate}
}

func hashFieldLines(field, value string) []string {
	return []string{
		fmt.Sprintf(hashFieldTemplate, field, len(value)),
		value,
	}
}

func parseValue(lengthParam string, data io.Reader) (string, error) {
	length, err := strconv.Atoi(lengthParam)
	if err != nil {
		return "", err
	}
	value := make([]byte, length, length)
	n, err := data.Read(value)
	if err != nil || n != length {
		return "", err
	}
	return string(value), nil
}

func parseTTL(ttlParam string) (uint64, error) {
	return strconv.ParseUint(ttlParam, 10, 0)
}

func newKeysCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile("^$"),
		run: func(params []string, data io.Reader) []string {
			keys := storage.Keys()
			response := []string{dataTemplate}
			for _, key := range keys {
				response = append(response, fmt.Sprintf(keyTemplate, key))
			}
			response = append(response, endTemplate)
			return response
		},
	}
}

func newTTLCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string, data io.Reader) []string {
			ttl, err := storage.TTL(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return []string{dataTemplate, fmt.Sprintf(ttlTemplate, ttl), endTemplate}
		},
	}
}

func newGetCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string, data io.Reader) []string {
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
		format:       regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, intPattern, intPattern)),
		requiredData: true,
		run: func(params []string, data io.Reader) []string {
			ttl, err := parseTTL(params[1])
			if err != nil {
				return errorResponse(err)
			}
			value, err := parseValue(params[2], data)
			if err != nil {
				return errorResponse(err)
			}
			if err := storage.Set(params[0], value, ttl); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newUpdCommand(storage storage.Storage) *command {
	return &command{
		format:       regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, intPattern)),
		requiredData: true,
		run: func(params []string, data io.Reader) []string {
			value, err := parseValue(params[1], data)
			if err != nil {
				return errorResponse(err)
			}
			if err := storage.Update(params[0], value); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newDelCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string, data io.Reader) []string {
			if err := storage.Delete(params[0]); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashCreateCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, intPattern)),
		run: func(params []string, data io.Reader) []string {
			ttl, err := parseTTL(params[1])
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
		run: func(params []string, data io.Reader) []string {
			hash, err := storage.HashGetAll(params[0])
			if err != nil {
				return errorResponse(err)
			}
			response := []string{dataTemplate}
			for field, value := range hash {
				response = append(response, hashFieldLines(field, value)...)
			}
			response = append(response, endTemplate)
			return response
		},
	}
}

func newHashGetCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, keyPattern)),
		run: func(params []string, data io.Reader) []string {
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
		format:       regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, keyPattern, intPattern)),
		requiredData: true,
		run: func(params []string, data io.Reader) []string {
			value, err := parseValue(params[2], data)
			if err != nil {
				return errorResponse(err)
			}
			if err := storage.HashSet(params[0], params[1], value); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newHashDelCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, keyPattern)),
		run: func(params []string, data io.Reader) []string {
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
		run: func(params []string, data io.Reader) []string {
			len, err := storage.HashLen(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return lengthResponse(len)
		},
	}
}

func newHashKeysCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string, data io.Reader) []string {
			keys, err := storage.HashKeys(params[0])
			if err != nil {
				return errorResponse(err)
			}
			response := []string{dataTemplate}
			for _, key := range keys {
				response = append(response, fmt.Sprintf(keyTemplate, key))
			}
			response = append(response, endTemplate)
			return response
		},
	}
}

func newListCreateCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, intPattern)),
		run: func(params []string, data io.Reader) []string {
			ttl, err := parseTTL(params[1])
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
		run: func(params []string, data io.Reader) []string {
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
		run: func(params []string, data io.Reader) []string {
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
		format:       regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, intPattern)),
		requiredData: true,
		run: func(params []string, data io.Reader) []string {
			value, err := parseValue(params[1], data)
			if err != nil {
				return errorResponse(err)
			}
			if err := storage.ListLeftPush(params[0], value); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newListRightPushCommand(storage storage.Storage) *command {
	return &command{
		format:       regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, intPattern)),
		requiredData: true,
		run: func(params []string, data io.Reader) []string {
			value, err := parseValue(params[1], data)
			if err != nil {
				return errorResponse(err)
			}
			if err := storage.ListRightPush(params[0], value); err != nil {
				return errorResponse(err)
			}
			return okResponse
		},
	}
}

func newListLenCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(params []string, data io.Reader) []string {
			len, err := storage.ListLen(params[0])
			if err != nil {
				return errorResponse(err)
			}
			return lengthResponse(len)
		},
	}
}

func newListRangeCommand(storage storage.Storage) *command {
	return &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, intPattern, intPattern)),
		run: func(params []string, data io.Reader) []string {
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
			response := []string{dataTemplate}
			for _, value := range values {
				response = append(response, valueLines(value)...)
			}
			response = append(response, endTemplate)
			return response
		},
	}
}

func newAuthCommand(htpasswdFile *htpasswd.HtpasswdFile, session *session) *command {
	return &command{
		allowGuest: true,
		format:     regexp.MustCompile("^([a-zA-Z0-9]+) (.+)$"),
		run: func(params []string, data io.Reader) []string {
			if htpasswdFile == nil || htpasswdFile.Validate(params[0], params[1]) {
				session.authorize()
				return okResponse
			}
			return errorResponse("INVALID CREDENTIALS")
		},
	}
}
