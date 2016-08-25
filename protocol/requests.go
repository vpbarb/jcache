package protocol

import (
	"errors"
	"fmt"
	"io"
	"regexp"
)

const (
	keyTemplate = "[a-zA-Z0-9_]+"
)

var (
	invalidRequestFormatError  = errors.New("Invalid request format")
	invalidUserFormatError     = errors.New("User is not valid")
	invalidPasswordFormatError = errors.New("Password is not valid")
	invalidKeyFormatError      = errors.New("Key is not valid")
	invalidFieldFormatError    = errors.New("Field is not valid")
	invalidValueLengthError    = errors.New("Value length is invalid")

	keyRegexp = regexp.MustCompile("^" + keyTemplate + "$")
)

type request struct {
	command string
}

func (r request) Command() string {
	return r.command
}

func (r *request) Decode(reader io.Reader) error {
	return readRequestEnd(reader)
}

func (r *request) Encode(writer io.Writer) (err error) {
	_, err = writer.Write([]byte(fmt.Sprintf("%s\r\n", r.command)))
	return
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
		return invalidUserFormatError
	}
	if !keyRegexp.MatchString(r.Password) {
		return invalidPasswordFormatError
	}
	return nil
}

func (r *authRequest) Decode(reader io.Reader) error {
	var user string
	var password string

	_, err := fmt.Fscanf(reader, "%s %s\r\n", &user, &password)
	if err != nil {
		return invalidRequestFormatError
	}

	r.User = user
	r.Password = password
	return nil
}

func (r *authRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s %s\r\n", r.command, r.User, r.Password)))
	return
}

type keyRequest struct {
	request
	Key string
}

func (r *keyRequest) validate() error {
	if keyRegexp.MatchString(r.Key) {
		return nil
	}
	return invalidKeyFormatError
}

func (r *keyRequest) Decode(reader io.Reader) error {
	var key string

	_, err := fmt.Fscanf(reader, "%s\r\n", &key)
	if err != nil {
		return invalidRequestFormatError
	}

	r.Key = key
	return nil
}

func (r *keyRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s\r\n", r.command, r.Key)))
	return
}

func newKeyRequest(command string) *keyRequest {
	return &keyRequest{request: newRequest(command)}
}

type keyTTLRequest struct {
	*keyRequest
	TTL uint64
}

func (r *keyTTLRequest) Decode(reader io.Reader) error {
	var key string
	var ttl uint64

	_, err := fmt.Fscanf(reader, "%s %d\r\n", &key, &ttl)
	if err != nil {
		return invalidRequestFormatError
	}

	r.Key = key
	r.TTL = ttl
	return nil
}

func (r *keyTTLRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s %d\r\n", r.command, r.Key, r.TTL)))
	return
}

type keyValueRequest struct {
	*keyRequest
	Value string
}

func newKeyTTLRequest(command string) *keyTTLRequest {
	return &keyTTLRequest{keyRequest: newKeyRequest(command)}
}

func (r *keyValueRequest) Decode(reader io.Reader) error {
	var key string
	var length int

	_, err := fmt.Fscanf(reader, "%s %d\r\n", &key, &length)
	if err != nil {
		return invalidRequestFormatError
	}

	value, err := readRequestValue(reader, length)
	if err != nil {
		return err
	}

	r.Key = key
	r.Value = string(value)
	return nil
}

func (r *keyValueRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s %d\r\n%s\r\n", r.command, r.Key, len(r.Value), r.Value)))
	return
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
		return invalidFieldFormatError
	}
	return nil
}

func (r *keyFieldRequest) Decode(reader io.Reader) error {
	var key string
	var field string

	_, err := fmt.Fscanf(reader, "%s %s\r\n", &key, &field)
	if err != nil {
		return invalidRequestFormatError
	}

	r.Key = key
	r.Field = field
	return nil
}

func (r *keyFieldRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s %s\r\n", r.command, r.Key, r.Field)))
	return
}

func newKeyFieldRequest(command string) *keyFieldRequest {
	return &keyFieldRequest{keyRequest: newKeyRequest(command)}
}

type keyFieldValueRequest struct {
	*keyFieldRequest
	Value string
}

func (r *keyFieldValueRequest) Decode(reader io.Reader) error {
	var key string
	var field string
	var length int

	_, err := fmt.Fscanf(reader, "%s %s %d\r\n", &key, &field, &length)
	if err != nil {
		return invalidRequestFormatError
	}

	value, err := readRequestValue(reader, length)
	if err != nil {
		return err
	}

	r.Key = key
	r.Field = field
	r.Value = string(value)
	return nil
}

func (r *keyFieldValueRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s %s %d\r\n%s\r\n", r.command, r.Key, r.Field, len(r.Value), r.Value)))
	return
}

func newKeyFieldValueRequest(command string) *keyFieldValueRequest {
	return &keyFieldValueRequest{keyFieldRequest: newKeyFieldRequest(command)}
}

type setRequest struct {
	*keyValueRequest
	TTL uint64
}

func (r *setRequest) Decode(reader io.Reader) error {
	var key string
	var ttl uint64
	var length int

	_, err := fmt.Fscanf(reader, "%s %d %d\r\n", &key, &ttl, &length)
	if err != nil {
		return invalidRequestFormatError
	}

	value, err := readRequestValue(reader, length)
	if err != nil {
		return err
	}

	r.Key = key
	r.TTL = ttl
	r.Value = string(value)
	return nil
}

func (r *setRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s %d %d\r\n%s\r\n", r.command, r.Key, r.TTL, len(r.Value), r.Value)))
	return
}

type listRangeRequest struct {
	*keyRequest
	Start int
	Stop  int
}

func (r *listRangeRequest) Decode(reader io.Reader) error {
	var key string
	var start, stop int

	_, err := fmt.Fscanf(reader, "%s %d %d\r\n", &key, &start, &stop)
	if err != nil {
		return invalidRequestFormatError
	}

	r.Key = key
	r.Start = start
	r.Stop = stop
	return nil
}

func (r *listRangeRequest) Encode(writer io.Writer) (err error) {
	if err := r.validate(); err != nil {
		return err
	}
	_, err = writer.Write([]byte(fmt.Sprintf("%s %s %d %d\r\n", r.command, r.Key, r.Start, r.Stop)))
	return
}

func readRequestValue(reader io.Reader, length int) ([]byte, error) {
	value := make([]byte, length, length)
	n, err := reader.Read(value)
	if err != nil {
		return nil, invalidRequestFormatError
	}
	if n != length {
		return nil, invalidValueLengthError
	}
	if err := readRequestEnd(reader); err != nil {
		return nil, err
	}
	return value, nil
}

func readRequestEnd(reader io.Reader) error {
	_, err := fmt.Fscanf(reader, "\r\n")
	if err != nil {
		return invalidRequestFormatError
	}
	return nil
}
