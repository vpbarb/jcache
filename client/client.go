package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/fatih/pool.v2"
)

const (
	errorPrefix = "-"
	countPrefix = "$"
	okResponse  = "+"
)

var (
	valueRegexp     = regexp.MustCompile(`^"(.*)"$`)
	hashFieldRegexp = regexp.MustCompile(`^([a-zA-Z0-9_]+):"(.*)"$`)
)

type Client struct {
	pool pool.Pool
}

func NewClient(addr, user, password string, timeout time.Duration, maxConnections int) (*Client, error) {
	factory := func() (net.Conn, error) {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return nil, fmt.Errorf("Cannot connect: %s", err)
		}

		response, err := call(conn, fmt.Sprintf("AUTH %s %s", user, password))
		if err != nil {
			return nil, err
		}

		switch response[0] {
		case okResponse:
			return conn, nil
		case "INVALID AUTH":
			return nil, fmt.Errorf("Invalid authentitication")
		default:
			return nil, fmt.Errorf("Invalid response")
		}
	}

	if pool, err := pool.NewChannelPool(0, maxConnections, factory); err == nil {
		return &Client{pool: pool}, nil
	} else {
		return nil, fmt.Errorf("Cannot create connection pool: %s", err)
	}
}

func (c *Client) Get(key string) (string, error) {
	conn, err := c.pool.Get()
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

func call(conn net.Conn, command string) ([]string, error) {
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
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

		if strings.HasPrefix(str, errorPrefix) {
			return nil, errors.New(strings.TrimPrefix(str, errorPrefix))
		}

		response = append(response, str)
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
