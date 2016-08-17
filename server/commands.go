package server

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	okResponse             = "OK"
	unknownCommandResponse = "UNKNOWN COMMAND"
	invalidFormatResponse  = "INVALID FORMAT"
	errorResponse          = "ERROR: %s"

	keyRegexp         = regexp.MustCompile(`(?i)^([a-z0-9_]+)$`)
	keyValueRegexp    = regexp.MustCompile(`(?i)^([a-z0-9_]+) "(.*)"$`)
	keyValueTTLRegexp = regexp.MustCompile(`(?i)^([a-z0-9_]+) "(.*)" ([a-z0-9-.]+)?$`)
)

func keysCommand(session *session, params string) string {
	keys := session.server.storage.Keys()
	return strings.Join(keys, "\n")
}

func ttlCommand(session *session, params string) string {
	matches := keyRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	ttl, err := session.server.storage.TTL(matches[1])
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return ttl.String()
}

func getCommand(session *session, params string) string {
	matches := keyRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	value, err := session.server.storage.Get(matches[1])
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return value
}

func setCommand(session *session, params string) string {
	matches := keyValueTTLRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	ttl, err := time.ParseDuration(matches[3])
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	if err := session.server.storage.Set(matches[1], matches[2], ttl); err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return okResponse
}

func delCommand(session *session, params string) string {
	matches := keyRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	if err := session.server.storage.Delete(matches[1]); err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return okResponse
}
