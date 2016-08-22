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

	keyRegexp = regexp.MustCompile(keyTemplate)
	intRegexp = regexp.MustCompile(intTemplate)

	requestHeaderRegexp          = regexp.MustCompile("^([A-Z]+)$")
	keyRequestHeaderRegexp       = regexp.MustCompile(fmt.Sprintf("^([A-Z]+) (%s)$", keyTemplate))
	keyIntRequestHeaderRegexp    = regexp.MustCompile(fmt.Sprintf("^([A-Z]+) (%s) (%s)$", keyTemplate, intTemplate))
	keyKeyRequestHeaderRegexp    = regexp.MustCompile(fmt.Sprintf("^([A-Z]+) (%s) (%s)$", keyTemplate, keyTemplate))
	keyKeyIntRequestHeaderRegexp = regexp.MustCompile(fmt.Sprintf("^([A-Z]+) (%s) (%s) (%s)$", keyTemplate, keyTemplate, intTemplate))
	keyIntIntRequestHeaderRegexp = regexp.MustCompile(fmt.Sprintf("^([A-Z]+) (%s) (%s) (%s)$", keyTemplate, intTemplate, intTemplate))
)

type request struct {
	command string
}

func (r request) checkCommand(command string) error {
	if command != r.command {
		return fmt.Errorf("Invalid command")
	}
	return nil
}

func (r *request) Decode(header []byte, data io.Reader) error {
	matches := requestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 2 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	return nil
}

func (r *request) Encode() ([]byte, error) {
	return []byte(fmt.Sprintf("%s\r\n", r.command)), nil
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
	matches := keyKeyRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 4 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	r.User = matches[2]
	r.Password = matches[3]
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
	if len(matches) < 3 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	r.Key = matches[2]
	return nil
}

func (r *keyRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s\r\n", r.command, r.Key)), nil
}

func newKeyRequest(command string) *keyRequest {
	return &keyRequest{request: newRequest(command)}
}

type keyTTLRequest struct {
	*keyRequest
	TTL uint64
}

func (r *keyTTLRequest) Decode(header []byte, data io.Reader) error {
	matches := keyIntRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 4 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	ttl, err := parseTTL(matches[3])
	if err != nil {
		return invalidCommandFormatError
	}
	r.Key = matches[2]
	r.TTL = ttl
	return nil
}

func (r *keyTTLRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %d\r\n", r.command, r.Key, r.TTL)), nil
}

type keyValueRequest struct {
	*keyRequest
	Value string
}

func (r *keyValueRequest) Decode(header []byte, data io.Reader) error {
	matches := keyIntRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 4 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	value, err := parseValue(matches[3], data)
	if err != nil {
		return invalidCommandFormatError
	}
	r.Key = matches[2]
	r.Value = value
	return nil
}

func (r *keyValueRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %d\r\n%s\r\n", r.command, r.Key, len(r.Value), r.Value)), nil
}

func newKeyValueRequest(command string) *keyValueRequest {
	return &keyValueRequest{keyRequest: newKeyRequest(command)}
}

type keyFieldRequest struct {
	*keyRequest
	Field string
}

func (r *keyFieldRequest) validate() error {
	if err := r.keyRequest.validate(); err != nil {
		return err
	}
	if !keyRegexp.MatchString(r.Field) {
		return fmt.Errorf("Field is not valid")
	}
	return nil
}

func (r *keyFieldRequest) Decode(header []byte, data io.Reader) error {
	matches := keyKeyRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 4 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	r.Key = matches[2]
	r.Field = matches[3]
	return nil
}

func (r *keyFieldRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %s\r\n", r.command, r.Key, r.Field)), nil
}

func newKeyFieldRequest(command string) *keyFieldRequest {
	return &keyFieldRequest{keyRequest: newKeyRequest(command)}
}

type keyFieldValueRequest struct {
	*keyFieldRequest
	Value string
}

func (r *keyFieldValueRequest) Decode(header []byte, data io.Reader) error {
	matches := keyKeyIntRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 5 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	value, err := parseValue(matches[4], data)
	if err != nil {
		return invalidCommandFormatError
	}
	r.Key = matches[2]
	r.Field = matches[3]
	r.Value = value
	return nil
}

func (r *keyFieldValueRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %s %d\r\n%s\r\n", r.command, r.Key, r.Field, len(r.Value), r.Value)), nil
}

type setRequest struct {
	*keyValueRequest
	TTL uint64
}

func (r *setRequest) Decode(header []byte, data io.Reader) error {
	matches := keyIntIntRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 4 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	ttl, err := parseTTL(matches[3])
	if err != nil {
		return invalidCommandFormatError
	}
	value, err := parseValue(matches[4], data)
	if err != nil {
		return invalidCommandFormatError
	}
	r.Key = matches[2]
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

type listRangeRequest struct {
	*keyRequest
	Start int
	Stop  int
}

func (r *listRangeRequest) Decode(header []byte, data io.Reader) error {
	matches := keyIntIntRequestHeaderRegexp.FindStringSubmatch(string(header))
	if len(matches) < 5 {
		return invalidCommandFormatError
	}
	if err := r.checkCommand(matches[1]); err != nil {
		return err
	}
	start, err := strconv.Atoi(matches[3])
	if err != nil {
		return invalidCommandFormatError
	}
	stop, err := strconv.Atoi(matches[4])
	if err != nil {
		return invalidCommandFormatError
	}
	r.Key = matches[2]
	r.Start = start
	r.Stop = stop
	return nil
}

func (r *listRangeRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %d %d\r\n", r.command, r.Key, r.Start, r.Stop)), nil
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
