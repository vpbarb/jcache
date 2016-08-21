package protocol

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type request struct {
	command string
}

func (r request) Command() string {
	return r.command
}

var (
	invalidCommandFormatError = fmt.Errorf("Invalid command format")

	keyRequestHeaderRegexp = regexp.MustCompile(fmt.Sprintf("^[A-Z]+ (%s)$", keyTemplate))
	setRequestHeaderRegexp = regexp.MustCompile(fmt.Sprintf("^[A-Z]+ (%s) (%s) (%s)$", keyTemplate, intTemplate, intTemplate))
)

type keyRequest struct {
	request
	Key string
}

func (r *keyRequest) validate() error {
	if keyRegexp.MatchString(r.Key) {
		return nil
	}
	return fmt.Errorf("Key is not valid")
}

func (r *keyRequest) Decode(header []byte, data io.Reader) error {
	matches := keyRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 2 {
		return invalidCommandFormatError
	}
	r.Key = matches[1]
	return nil
}

func (r *keyRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s\r\n", r.command, r.Key)), nil
}

func newKeyRequest(command, key string) *keyRequest {
	return &keyRequest{request{command: command}, key}
}

func NewGetRequest(key string) *keyRequest {
	return newKeyRequest("GET", key)
}

func NewSetRequest(key, value string, ttl uint64) *setRequest {
	return &setRequest{
		keyRequest: newKeyRequest("SET", key),
		Value:      value,
		TTL:        ttl,
	}
}

func NewDelRequest(key string) *keyRequest {
	return newKeyRequest("DEL", key)
}

type setRequest struct {
	*keyRequest
	Value string
	TTL   uint64
}

func (r *setRequest) Decode(header []byte, data io.Reader) error {
	matches := setRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 3 {
		return invalidCommandFormatError
	}
	ttl, err := parseTTL(matches[2])
	if err != nil {
		return invalidCommandFormatError
	}
	value, err := parseValue(matches[3], data)
	if err != nil {
		return invalidCommandFormatError
	}
	r.Key = matches[1]
	r.TTL = ttl
	r.Value = value
	return nil
}

func (r *setRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %d %d\r\n%s\r\n", r.command, r.Key, r.TTL, len(r.Value), r.Value)), nil
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
