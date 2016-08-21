package client

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func transfer(w io.Writer, r io.Reader, request string, dataFormatter dataFormatFunc) error {
	wb := bufio.NewWriter(w)
	wb.WriteString(request)
	if err := wb.Flush(); err != nil {
		return fmt.Errorf("Cannot write to connection: %s", err)
	}

	rb := bufio.NewReader(r)
	line, _, err := rb.ReadLine()
	if err != nil {
		return fmt.Errorf("Cannot read from connection: %s", err)
	}
	str := string(line)
	switch {
	case strings.HasPrefix(str, "ERROR"):
		return fmt.Errorf("Response error: %s", strings.TrimPrefix(str, "ERROR "))
	case str == "OK":
		return nil
	case str == "DATA":
		return dataFormatter(rb)
	default:
		return fmt.Errorf("Invalid response format")
	}
}
