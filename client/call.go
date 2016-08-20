package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	okPrefix    = "+"
	errorPrefix = "-"
)

func call(w io.Writer, r io.Reader, command string) ([]string, error) {
	wb := bufio.NewWriter(w)
	rb := bufio.NewReader(r)
	wb.WriteString(command + "\r\n")
	if err := wb.Flush(); err != nil {
		return nil, fmt.Errorf("Cannot write to connection: %s", err)
	}

	var response []string
	for {
		line, _, err := rb.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("Cannot read from connection: %s", err)
		}
		str := string(line)
		switch {
		case str == "":
			// Ignore empty line and read next one
			continue
		case strings.HasPrefix(str, okPrefix):
			return response, nil
		case strings.HasPrefix(str, errorPrefix):
			return nil, errors.New(strings.TrimPrefix(str, errorPrefix))
		default:
			response = append(response, str)
		}
	}

	return response, nil
}
