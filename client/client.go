package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/fatih/pool.v2"
)

const (
	okPrefix    = "+"
	errorPrefix = "-"
	countPrefix = "$"
)

var (
	valueRegexp     = regexp.MustCompile(`^"(.*)"$`)
	hashFieldRegexp = regexp.MustCompile(`^([a-zA-Z0-9_]+):"(.*)"$`)
)

type Client struct {
	connPool pool.Pool
}

func NewClient(addr, user, password string, timeout time.Duration, maxConnections int) (*Client, error) {
	factory := func() (net.Conn, error) {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return nil, fmt.Errorf("Cannot connect: %s", err)
		}

		_, err = call(conn, fmt.Sprintf("AUTH %s %s", user, password))
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("Cannot authentiticate: %s", err)
		}

		return conn, nil
	}

	if connPool, err := pool.NewChannelPool(0, maxConnections, factory); err == nil {
		return &Client{connPool: connPool}, nil
	} else {
		return nil, fmt.Errorf("Cannot create connection pool: %s", err)
	}
}

func (c *Client) Get(key string) (string, error) {
	conn, err := c.connPool.Get()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	response, err := call(conn, fmt.Sprintf("GET %s", key))
	if err != nil {
		return "", err
	}

	return parseValue(response[0]), nil
}

func call(rw io.ReadWriter, command string) ([]string, error) {
	w := bufio.NewWriter(rw)
	r := bufio.NewReader(rw)
	w.WriteString(command + "\r\n")
	if err := w.Flush(); err != nil {
		return nil, fmt.Errorf("Cannot write to connection: %s", err)
	}

	var response []string
	var i, count int
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("Cannot read from connection: %s", err)
		}
		str := string(line)
		if strings.HasPrefix(str, countPrefix) {
			count, err = strconv.Atoi(strings.TrimPrefix(str, countPrefix))
			if err != nil {
				return nil, fmt.Errorf("Invalid rows count")
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

func parseValue(str string) string {
	if matches := valueRegexp.FindStringSubmatch(str); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func parseHashField(str string) (string, string) {
	if matches := hashFieldRegexp.FindStringSubmatch(str); len(matches) > 2 {
		return matches[1], matches[2]
	}
	return "", ""
}
