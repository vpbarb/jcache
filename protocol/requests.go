package protocol

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
)

const (
	keyTemplate = "[a-zA-Z0-9_]+"
	intTemplate = "[0-9_]+"
)

var (
	invalidCommandFormatError = fmt.Errorf("Invalid command format")

	keyRegexp               = regexp.MustCompile(keyTemplate)
	ttlRegexp               = regexp.MustCompile(intTemplate)
	lenRegexp               = regexp.MustCompile(intTemplate)
	keyRequestHeaderRegexp  = regexp.MustCompile(fmt.Sprintf("^[A-Z]+ (%s)$", keyTemplate))
	authRequestHeaderRegexp = regexp.MustCompile(fmt.Sprintf("^[A-Z]+ (%s) (%s)$", keyTemplate, keyTemplate))
	setRequestHeaderRegexp  = regexp.MustCompile(fmt.Sprintf("^[A-Z]+ (%s) (%s) (%s)$", keyTemplate, intTemplate, intTemplate))
)

type request struct {
	command string
}

func newRequest(command string) request {
	return request{command: command}
}

type authRequest struct {
	request
	User     string
	Password string
}

func (r *authRequest) validate() error {
	if !keyRegexp.MatchString(r.User) {
		return fmt.Errorf("User is not valid")
	}
	if !keyRegexp.MatchString(r.Password) {
		return fmt.Errorf("Password is not valid")
	}
	return nil
}

func (r *authRequest) Decode(header []byte, data io.Reader) error {
	matches := authRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 3 {
		return invalidCommandFormatError
	}
	r.User = matches[1]
	r.Password = matches[2]
	return nil
}

func (r *authRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %s\r\n", r.command, r.User, r.Password)), nil
}

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
	return &keyRequest{request: newRequest(command), Key: key}
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
