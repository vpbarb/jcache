package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	dataKeyword = "DATA"
	endKeyword  = "END"
	okKeyword   = "OK"
	errorPrefix = "ERROR "
)

var (
	invalidResponseFormatError = fmt.Errorf("Invalid response format")
	invalidDataFormatError     = fmt.Errorf("Invalid data format")
)

type response struct {
	Error error
}

func (r *response) encodeResponse(response []byte) ([]byte, error) {
	if r.Error != nil {
		return []byte(fmt.Sprintf("%s%s\r\n", errorPrefix, r.Error)), nil
	}
	return response, nil
}

func (r *response) decodeHeader(buf *bufio.Reader) ([]byte, error) {
	header, _, err := buf.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("Cannot read header: %s", err)
	}
	if strings.HasPrefix(string(header), errorPrefix) {
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(string(header), errorPrefix))
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
	return r.encodeResponse([]byte(okKeyword + "\r\n"))
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
	if string(header) == okKeyword {
		return nil
	}
	return invalidResponseFormatError
}

type dataResponse struct {
	*response
}

func newDataResponse() dataResponse {
	return dataResponse{response: &response{}}
}

func (r dataResponse) encodeData(data []byte) ([]byte, error) {
	response := []byte(dataKeyword + "\r\n")
	response = append(response, data...)
	response = append(response, []byte(endKeyword+"\r\n")...)
	return r.encodeResponse(response)
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
	dataResponse
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
	return r.encodeData(data)
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
	if string(header) == dataKeyword {
		var keys []string
		for {
			header, err := r.decodeHeader(buf)
			if err != nil {
				return err
			}
			if string(header) == endKeyword {
				r.Keys = keys
				return nil
			}
			var key string
			_, err = fmt.Sscanf(string(header), "KEY %s", &key)
			if err != nil {
				return invalidResponseFormatError
			}
			keys = append(keys, key)
		}
	}
	return invalidResponseFormatError
}

type valuesResponse struct {
	dataResponse
	Values []string
}

func (r *valuesResponse) Encode() ([]byte, error) {
	var data []byte
	for _, value := range r.Values {
		data = append(data, []byte(fmt.Sprintf("VALUE %d\r\n%s\r\n", len(value), value))...)
	}
	return r.encodeData(data)
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
	if string(header) == dataKeyword {
		var values []string
		for {
			header, err := r.decodeHeader(buf)
			if err != nil {
				return err
			}
			if string(header) == endKeyword {
				r.Values = values
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
			values = append(values, string(value))
		}
	}
	return invalidResponseFormatError
}

type fieldsResponse struct {
	dataResponse
	Fields map[string]string
}

func (r *fieldsResponse) Encode() ([]byte, error) {
	var data []byte
	for field, value := range r.Fields {
		data = append(data, []byte(fmt.Sprintf("FIELD %s %d\r\n%s\r\n", field, len(value), value))...)
	}
	return r.encodeData(data)
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
	if string(header) == dataKeyword {
		fields := make(map[string]string)
		for {
			header, err := r.decodeHeader(buf)
			if err != nil {
				return err
			}
			if string(header) == endKeyword {
				r.Fields = fields
				return nil
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
	}
	return invalidResponseFormatError
}

func readResponseValue(buf *bufio.Reader, length int) (string, error) {
	value := make([]byte, length, length)
	n, err := buf.Read(value)
	if err != nil {
		return "", err
	}
	if n != length {
		return "", fmt.Errorf("Value length is invalid")
	}
	rest, _, err := buf.ReadLine()
	if len(rest) > 0 || err != nil {
		return "", fmt.Errorf("Value length is invalid")
	}
	return string(value), nil
}
