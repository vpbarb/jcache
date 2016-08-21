package protocol

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	invalidResponseFormatError = fmt.Errorf("Invalid response format")
	invalidDataFormatError     = fmt.Errorf("Invalid data format")

	valueHeaderRegexp = regexp.MustCompile("^VALUE ([0-9]+)$")
)

type response struct {
	Error error
}

func (r response) Encode() ([]byte, error) {
	return r.encodeResponse([]byte("OK\r\n"))
}

func (r response) encodeResponse(response []byte) ([]byte, error) {
	if r.Error != nil {
		return []byte(fmt.Sprintf("ERROR %s\r\n", r.Error)), nil
	}
	return response, nil
}

func (r response) encodeData(data []byte) ([]byte, error) {
	response := []byte("DATA\r\n")
	response = append(response, data...)
	response = append(response, []byte("END\r\n")...)
	return r.encodeResponse(response)
}

func (r response) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, "ERROR"):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, "ERROR "))
		return nil
	case str == "OK":
		return nil
		return nil
	}
	return invalidResponseFormatError
}

type valueResponse struct {
	response
	Value string
}

func (r *valueResponse) Encode() ([]byte, error) {
	return r.encodeData([]byte(fmt.Sprintf("VALUE %d\r\n%s\r\n", len(r.Value), r.Value)))
}

func (r *valueResponse) Decode(header []byte, data io.Reader) error {
	str := string(header)
	switch {
	case strings.HasPrefix(str, "ERROR"):
		r.Error = fmt.Errorf("Response error: %s", strings.TrimPrefix(str, "ERROR "))
		return nil
	case str == "DATA":
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
			return nil
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
	if str == "END" {
		return "", true, nil
	}
	matches := valueHeaderRegexp.FindStringSubmatch(str)
	if len(matches) < 1 {
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
