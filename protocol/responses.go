package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	invalidResponseFormatError = errors.New("Invalid response format")
)

type response struct {
	Error error
}

func (r *response) encodeResponse(response []byte) ([]byte, error) {
	if r.Error != nil {
		return []byte(fmt.Sprintf("ERROR %s\r\n", r.Error)), nil
	}
	return response, nil
}

func (r *response) decodeHeader(buf *bufio.Reader) ([]byte, error) {
	header, _, err := buf.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("Cannot read header: %s", err)
	}
	if strings.HasPrefix(string(header), "ERROR") {
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(string(header), "ERROR "))
	}
	return header, nil
}

type okResponse struct {
	*response
}

func newOkResponse() *okResponse {
	return &okResponse{response: &response{}}
}

func (r *okResponse) Encode() ([]byte, error) {
	return r.encodeResponse([]byte("OK\r\n"))
}

func (r *okResponse) Decode(data io.Reader) error {
	buf := bufio.NewReader(data)
	header, err := r.decodeHeader(buf)
	if err != nil {
		return err
	}
	if r.Error != nil {
		return nil
	}
	if string(header) == "OK" {
		return nil
	}
	return invalidResponseFormatError
}

type countResponse struct {
	*response
}

func newDataResponse() countResponse {
	return countResponse{response: &response{}}
}

func (r countResponse) encodeData(data []byte, count int) ([]byte, error) {
	response := []byte(fmt.Sprintf("COUNT %d\r\n", count))
	response = append(response, data...)
	return r.encodeResponse(response)
}

func (r countResponse) decodeCount(header []byte) (int, error) {
	var count int
	_, err := fmt.Sscanf(string(header), "COUNT %d", &count)
	if err != nil {
		return 0, invalidResponseFormatError
	}
	return count, nil
}

type lenResponse struct {
	*response
	Len int
}

func (r *lenResponse) Encode() ([]byte, error) {
	return r.encodeResponse([]byte(fmt.Sprintf("LEN %d\r\n", r.Len)))
}

func (r *lenResponse) Decode(data io.Reader) error {
	buf := bufio.NewReader(data)
	header, err := r.decodeHeader(buf)
	if err != nil {
		return err
	}
	if r.Error != nil {
		return nil
	}

	var length int
	_, err = fmt.Sscanf(string(header), "LEN %d", &length)
	if err != nil {
		return invalidResponseFormatError
	}
	r.Len = length
	return nil
}

type valueResponse struct {
	*response
	Value string
}

func newValueResponse() *valueResponse {
	return &valueResponse{response: &response{}}
}

func (r *valueResponse) Encode() ([]byte, error) {
	return r.encodeResponse([]byte(fmt.Sprintf("VALUE %d\r\n%s\r\n", len(r.Value), r.Value)))
}

func (r *valueResponse) Decode(data io.Reader) error {
	buf := bufio.NewReader(data)
	header, err := r.decodeHeader(buf)
	if err != nil {
		return err
	}
	if r.Error != nil {
		return nil
	}
	var length int
	_, err = fmt.Sscanf(string(header), "VALUE %d", &length)
	if err != nil {
		return invalidResponseFormatError
	}
	value, err := readResponseValue(buf, length)
	if err != nil {
		return err
	}
	r.Value = string(value)
	return nil
}

type keysResponse struct {
	countResponse
	Keys []string
}

func (r *keysResponse) Encode() ([]byte, error) {
	var data []byte
	for _, key := range r.Keys {
		if !keyRegexp.MatchString(key) {
			return nil, fmt.Errorf("Invalid key: %s", key)
		}
		data = append(data, []byte(fmt.Sprintf("KEY %s\r\n", key))...)
	}
	return r.encodeData(data, len(r.Keys))
}

func (r *keysResponse) Decode(data io.Reader) error {
	buf := bufio.NewReader(data)
	header, err := r.decodeHeader(buf)
	if err != nil {
		return err
	}
	if r.Error != nil {
		return nil
	}
	count, err := r.decodeCount(header)
	if err != nil {
		return err
	}
	var keys []string
	for i := 0; i < count; i++ {
		header, _, err := buf.ReadLine()
		if err != nil {
			return err
		}
		var key string
		_, err = fmt.Sscanf(string(header), "KEY %s", &key)
		if err != nil {
			return invalidResponseFormatError
		}
		keys = append(keys, key)
	}
	r.Keys = keys
	return nil
}

type valuesResponse struct {
	countResponse
	Values []string
}

func (r *valuesResponse) Encode() ([]byte, error) {
	var data []byte
	for _, value := range r.Values {
		data = append(data, []byte(fmt.Sprintf("VALUE %d\r\n%s\r\n", len(value), value))...)
	}
	return r.encodeData(data, len(r.Values))
}

func (r *valuesResponse) Decode(data io.Reader) error {
	buf := bufio.NewReader(data)
	header, err := r.decodeHeader(buf)
	if err != nil {
		return err
	}
	if r.Error != nil {
		return nil
	}
	count, err := r.decodeCount(header)
	if err != nil {
		return err
	}
	var values []string
	for i := 0; i < count; i++ {
		header, _, err := buf.ReadLine()
		if err != nil {
			return err
		}
		var length int
		_, err = fmt.Sscanf(string(header), "VALUE %d", &length)
		if err != nil {
			return invalidResponseFormatError
		}
		value, err := readResponseValue(buf, length)
		if err != nil {
			return err
		}
		values = append(values, string(value))
	}
	r.Values = values
	return nil
}

type fieldsResponse struct {
	countResponse
	Fields map[string]string
}

func (r *fieldsResponse) Encode() ([]byte, error) {
	var data []byte
	for field, value := range r.Fields {
		data = append(data, []byte(fmt.Sprintf("FIELD %s %d\r\n%s\r\n", field, len(value), value))...)
	}
	return r.encodeData(data, len(r.Fields))
}

func (r *fieldsResponse) Decode(data io.Reader) error {
	buf := bufio.NewReader(data)
	header, err := r.decodeHeader(buf)
	if err != nil {
		return err
	}
	if r.Error != nil {
		return nil
	}
	count, err := r.decodeCount(header)
	if err != nil {
		return err
	}
	fields := make(map[string]string)
	for i := 0; i < count; i++ {
		header, _, err := buf.ReadLine()
		if err != nil {
			return err
		}
		var field string
		var length int
		_, err = fmt.Sscanf(string(header), "FIELD %s %d", &field, &length)
		if err != nil {
			return invalidResponseFormatError
		}
		value, err := readResponseValue(buf, length)
		if err != nil {
			return err
		}
		fields[field] = string(value)
	}
	r.Fields = fields
	return nil
}

func readResponseValue(buf *bufio.Reader, length int) (string, error) {
	value := make([]byte, length, length)
	n, err := buf.Read(value)
	if err != nil {
		return "", err
	}
	if n != length {
		return "", invalidValueLengthError
	}
	rest, _, err := buf.ReadLine()
	if len(rest) > 0 || err != nil {
		return "", invalidValueLengthError
	}
	return string(value), nil
}
