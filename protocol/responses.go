package protocol

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
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

	valueHeaderRegexp = regexp.MustCompile("^VALUE ([0-9]+)$")
	fieldHeaderRegexp = regexp.MustCompile("^FIELD ([a-zA-Z0-9_]+) ([0-9]+)$")
	keyHeaderRegexp   = regexp.MustCompile("^KEY ([a-zA-Z0-9_]+)$")
	lenHeaderRegexp   = regexp.MustCompile("^LEN ([0-9]+)$")
)

type response struct {
	Error error
}

func (r response) encodeResponse(response []byte) ([]byte, error) {
	if r.Error != nil {
		return []byte(fmt.Sprintf("%s%s\r\n", errorPrefix, r.Error)), nil
	}
	return response, nil
}

type okResponse struct {
	response
}

func (r *okResponse) Encode() ([]byte, error) {
	return r.encodeResponse([]byte(okKeyword + "\r\n"))
}

func (r *okResponse) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, errorPrefix):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, errorPrefix))
		return nil
	case str == okKeyword:
		return nil
		return nil
	}
	return invalidResponseFormatError
}

type dataResponse struct {
	response
}

func (r dataResponse) encodeData(data []byte) ([]byte, error) {
	response := []byte(dataKeyword + "\r\n")
	response = append(response, data...)
	response = append(response, []byte(endKeyword+"\r\n")...)
	return r.encodeResponse(response)
}

type lenResponse struct {
	dataResponse
	Len int
}

func (r *lenResponse) Encode() ([]byte, error) {
	return r.encodeData([]byte(fmt.Sprintf("LEN %d\r\n", r.Len)))
}

func (r *lenResponse) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, errorPrefix):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, errorPrefix))
		return nil
	case str == dataKeyword:
		buf := bufio.NewReader(data)
		for {
			length, isEnd, err := readLen(buf)
			if err != nil {
				return err
			}
			if isEnd {
				return nil
			}
			r.Len = length
		}
	}
	return invalidResponseFormatError
}

type valueResponse struct {
	dataResponse
	Value string
}

func (r *valueResponse) Encode() ([]byte, error) {
	return r.encodeData([]byte(fmt.Sprintf("VALUE %d\r\n%s\r\n", len(r.Value), r.Value)))
}

func (r *valueResponse) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, errorPrefix):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, errorPrefix))
		return nil
	case str == dataKeyword:
		buf := bufio.NewReader(data)
		for {
			value, isEnd, err := readValue(buf)
			if err != nil {
				return err
			}
			if isEnd {
				return nil
			}
			r.Value = value
		}
	}
	return invalidResponseFormatError
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

func (r *keysResponse) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, errorPrefix):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, errorPrefix))
		return nil
	case str == dataKeyword:
		buf := bufio.NewReader(data)
		var keys []string
		for {
			key, isEnd, err := readKey(buf)
			if err != nil {
				return err
			}
			if isEnd {
				r.Keys = keys
				return nil
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

func (r *valuesResponse) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, errorPrefix):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, errorPrefix))
		return nil
	case str == dataKeyword:
		buf := bufio.NewReader(data)
		var values []string
		for {
			value, isEnd, err := readValue(buf)
			if err != nil {
				return err
			}
			if isEnd {
				r.Values = values
				return nil
			}
			values = append(values, value)
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

func (r *fieldsResponse) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, errorPrefix):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, errorPrefix))
		return nil
	case str == dataKeyword:
		buf := bufio.NewReader(data)
		fields := make(map[string]string)
		for {
			field, value, isEnd, err := readField(buf)
			if err != nil {
				return err
			}
			if isEnd {
				r.Fields = fields
				return nil
			}
			fields[field] = value
		}
	}
	return invalidResponseFormatError
}

func readValue(buf *bufio.Reader) (string, bool, error) {
	line, _, err := buf.ReadLine()
	if err != nil {
		return "", false, err
	}
	str := string(line)
	if str == endKeyword {
		return "", true, nil
	}
	matches := valueHeaderRegexp.FindStringSubmatch(str)
	if len(matches) < 2 {
		return "", false, invalidDataFormatError
	}
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", false, err
	}
	data := make([]byte, length, length)
	n, err := buf.Read(data)
	if err != nil || n != length {
		return "", false, err
	}
	value := string(data)
	buf.ReadLine()
	return value, false, nil
}

func readField(buf *bufio.Reader) (string, string, bool, error) {
	line, _, err := buf.ReadLine()
	if err != nil {
		return "", "", false, err
	}
	str := string(line)
	if str == endKeyword {
		return "", "", true, nil
	}
	matches := fieldHeaderRegexp.FindStringSubmatch(str)
	if len(matches) < 3 {
		return "", "", false, invalidDataFormatError
	}
	length, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", "", false, err
	}
	data := make([]byte, length, length)
	n, err := buf.Read(data)
	if err != nil || n != length {
		return "", "", false, err
	}
	value := string(data)
	buf.ReadLine()
	return matches[1], value, false, nil
}

func readKey(buf *bufio.Reader) (string, bool, error) {
	line, _, err := buf.ReadLine()
	if err != nil {
		return "", false, err
	}
	str := string(line)
	if str == endKeyword {
		return "", true, nil
	}
	matches := keyHeaderRegexp.FindStringSubmatch(str)
	if len(matches) < 2 {
		return "", false, invalidDataFormatError
	}
	key := matches[1]
	return key, false, nil
}

func readLen(buf *bufio.Reader) (int, bool, error) {
	line, _, err := buf.ReadLine()
	if err != nil {
		return 0, false, err
	}
	str := string(line)
	if str == endKeyword {
		return 0, true, nil
	}
	matches := lenHeaderRegexp.FindStringSubmatch(str)
	if len(matches) < 2 {
		return 0, false, invalidDataFormatError
	}
	length, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, false, err
	}
	return length, false, nil
}
