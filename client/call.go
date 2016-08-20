package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	okPrefix    = "+"
	errorPrefix = "-"
	countPrefix = "$"
)

func call(w io.Writer, r io.Reader, command string) ([]string, error) {
	wb := bufio.NewWriter(w)
	rb := bufio.NewReader(r)
	wb.WriteString(command + "\r\n")
	if err := wb.Flush(); err != nil {
		return nil, fmt.Errorf("Cannot write to connection: %s", err)
	}

	var response []string
	var i, count int
	for {
		line, _, err := rb.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("Cannot read from connection: %s", err)
		}
		str := string(line)
		if strings.HasPrefix(str, countPrefix) {
			count, err = strconv.Atoi(strings.TrimPrefix(str, countPrefix))
			if err != nil {
				return nil, errors.New("Invalid response rows count")
			}
			continue
		}
		if count == 0 {
			// Count can't be zero, look for count in the next line
			continue
		}

		switch string(str[0]) {
		case okPrefix:
			return nil, nil
		case errorPrefix:
			return nil, errors.New(strings.TrimPrefix(str, errorPrefix))
		default:
			response = append(response, str)
		}

		i++
		if i >= count {
			break
		}
	}

	return response, nil
}
