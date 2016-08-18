package server

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	errorResponse          = "ERROR: %s"
	valueResponse          = `"%s"`
	hashElementResponse    = `%s: "%s"`
)

type command struct {
	format     *regexp.Regexp
	run        func(session *session, params []string) string
	allowGuest bool
}

var (
	keysCommand = &command{
		format: regexp.MustCompile("^$"),
		run: func(session *session, params []string) string {
			keys := session.server.storage.Keys()
			return strings.Join(keys, "\n")
		},
	}

	ttlCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			ttl, err := session.server.storage.TTL(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return ttl.String()
		},
	}

	getCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.Get(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}

	setCommand = &command{
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

	updCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.Update(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}

	delCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.Delete(params[0]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}

	hashCreateCommand = &command{
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

	hashGetAllCommand = &command{
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

	hashGetCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.HashGet(params[0], params[1])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}

	hashSetCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s %s$", keyPattern, fieldPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.HashSet(params[0], params[1], params[2]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}

	hashDelCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, fieldPattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.HashDelete(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}

	hashLenCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			len, err := session.server.storage.HashLen(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf("%d", len)
		},
	}

	hashKeysCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			keys, err := session.server.storage.HashKeys(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return strings.Join(keys, "\n")
		},
	}

	listCreateCommand = &command{
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

	listLeftPopCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.ListLeftPop(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}

	listRightPopCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s$", keyPattern)),
		run: func(session *session, params []string) string {
			value, err := session.server.storage.ListRightPop(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf(valueResponse, value)
		},
	}

	listLeftPushCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.ListLeftPush(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}

	listRightPushCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^%s %s$", keyPattern, valuePattern)),
		run: func(session *session, params []string) string {
			if err := session.server.storage.ListRightPush(params[0], params[1]); err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return okResponse
		},
	}

	listLenCommand = &command{
		format: regexp.MustCompile(fmt.Sprintf("^LLEN %s$", keyPattern)),
		run: func(session *session, params []string) string {
			len, err := session.server.storage.ListLen(params[0])
			if err != nil {
				return fmt.Sprintf(errorResponse, err)
			}
			return fmt.Sprintf("%d", len)
		},
	}

	listRangeCommand = &command{
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
)
