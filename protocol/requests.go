package protocol

import (
	"fmt"
	"io"
	"regexp"
)

const (
	keyTemplate = "[a-zA-Z0-9_]+"
)

var (
	invalidCommandFormatError = fmt.Errorf("Invalid command format")

	keyRegexp = regexp.MustCompile("^" + keyTemplate + "$")
)

type request struct {
	command string
}

func (r request) Command() string {
	return r.command
}

func (r *request) Decode(data io.Reader) error {
	err := readRequestEnd(data)
	if err != nil {
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

func (r *authRequest) Decode(data io.Reader) error {
	var user string
	var password string

	_, err := fmt.Fscanf(data, "%s %s\r\n", &user, &password)
	if err != nil {
		return invalidCommandFormatError
	}

	r.User = user
	r.Password = password
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

func (r *keyRequest) Decode(data io.Reader) error {
	var key string

	_, err := fmt.Fscanf(data, "%s\r\n", &key)
	if err != nil {
		return invalidCommandFormatError
	}

	r.Key = key
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

func (r *keyTTLRequest) Decode(data io.Reader) error {
	var key string
	var ttl uint64

	_, err := fmt.Fscanf(data, "%s %d\r\n", &key, &ttl)
	if err != nil {
		return invalidCommandFormatError
	}

	r.Key = key
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

func newKeyTTLRequest(command string) *keyTTLRequest {
	return &keyTTLRequest{keyRequest: newKeyRequest(command)}
}

func (r *keyValueRequest) Decode(data io.Reader) error {
	var key string
	var length int

	_, err := fmt.Fscanf(data, "%s %d\r\n", &key, &length)
	if err != nil {
		return invalidCommandFormatError
	}

	value, err := readRequestValue(data, length)
	if err != nil {
		return err
	}

	r.Key = key
	r.Value = string(value)
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

func (r *keyFieldRequest) Decode(data io.Reader) error {
	var key string
	var field string

	_, err := fmt.Fscanf(data, "%s %s\r\n", &key, &field)
	if err != nil {
		return invalidCommandFormatError
	}

	r.Key = key
	r.Field = field
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

func (r *keyFieldValueRequest) Decode(data io.Reader) error {
	var key string
	var field string
	var length int

	_, err := fmt.Fscanf(data, "%s %s %d\r\n", &key, &field, &length)
	if err != nil {
		return invalidCommandFormatError
	}

	value, err := readRequestValue(data, length)
	if err != nil {
		return err
	}

	r.Key = key
	r.Field = field
	r.Value = string(value)
	return nil
}

func (r *keyFieldValueRequest) Encode() ([]byte, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s %s %s %d\r\n%s\r\n", r.command, r.Key, r.Field, len(r.Value), r.Value)), nil
}

func newKeyFieldValueRequest(command string) *keyFieldValueRequest {
	return &keyFieldValueRequest{keyFieldRequest: newKeyFieldRequest(command)}
}

type setRequest struct {
	*keyValueRequest
	TTL uint64
}

func (r *setRequest) Decode(data io.Reader) error {
	var key string
	var ttl uint64
	var length int

	_, err := fmt.Fscanf(data, "%s %d %d\r\n", &key, &ttl, &length)
	if err != nil {
		return invalidCommandFormatError
	}

	value, err := readRequestValue(data, length)
	if err != nil {
		return err
	}

	r.Key = key
	r.TTL = ttl
	r.Value = string(value)
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

func (r *listRangeRequest) Decode(data io.Reader) error {
	var key string
	var start, stop int

	_, err := fmt.Fscanf(data, "%s %d %d\r\n", &key, &start, &stop)
	if err != nil {
		return invalidCommandFormatError
	}

	r.Key = key
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

func readRequestValue(data io.Reader, length int) ([]byte, error) {
	value := make([]byte, length, length)
	n, err := data.Read(value)
	if err != nil {
		return nil, invalidCommandFormatError
	}
	if n != length {
		return nil, fmt.Errorf("Value length is invalid")
	}
	if err := readRequestEnd(data); err != nil {
		return nil, err
	}
	return value, nil
}

func readRequestEnd(data io.Reader) error {
	_, err := fmt.Fscanf(data, "\r\n")
	if err != nil {
		return invalidCommandFormatError
	}
	return nil
}
