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
	valueResponse          = `"%s"`
	hashElementResponse    = `%s: "%s"`

	keyRegexp           = regexp.MustCompile(`(?i)^([a-z0-9_]+)$`)
	keyValueRegexp      = regexp.MustCompile(`(?i)^([a-z0-9_]+) "(.*)"$`)
	keyValueTTLRegexp   = regexp.MustCompile(`(?i)^([a-z0-9_]+) "(.*)" ([a-z0-9-.]+)$`)
	keyTTLRegexp        = regexp.MustCompile(`(?i)^([a-z0-9_]+) ([a-z0-9-.]+)$`)
	keyFieldRegexp      = regexp.MustCompile(`(?i)^([a-z0-9_]+) ([a-z0-9_]+)$`)
	keyFieldValueRegexp = regexp.MustCompile(`(?i)^([a-z0-9_]+) ([a-z0-9_]+) "(.*)"$`)
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
	return fmt.Sprintf(valueResponse, value)
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

func hashCreateCommand(session *session, params string) string {
	matches := keyTTLRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	ttl, err := time.ParseDuration(matches[2])
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	err = session.server.storage.HashCreate(matches[1], ttl)
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return okResponse
}

func hashGetAllCommand(session *session, params string) string {
	matches := keyRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	hash, err := session.server.storage.HashGetAll(matches[1])
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	var response []string
	for key, value := range hash {
		response = append(response, fmt.Sprintf(hashElementResponse, key, value))
	}

	return strings.Join(response, "\n")
}

func hashGetCommand(session *session, params string) string {
	matches := keyFieldRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	value, err := session.server.storage.HashGet(matches[1], matches[2])
	if err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return fmt.Sprintf(valueResponse, value)
}

func hashSetCommand(session *session, params string) string {
	matches := keyFieldValueRegexp.FindStringSubmatch(params)
	if len(matches) == 0 {
		return invalidFormatResponse
	}
	if err := session.server.storage.HashSet(matches[1], matches[2], matches[3]); err != nil {
		return fmt.Sprintf(errorResponse, err)
	}
	return okResponse
}
