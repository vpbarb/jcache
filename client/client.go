package client

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"gopkg.in/fatih/pool.v2"
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

		_, err = callConn(conn, "AUTH %s %s", user, password)
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

	response, err := callConn(conn, "GET %s", key)
	if err != nil {
		return "", err
	}

	return parseValue(response[0]), nil
}

func callConn(conn net.Conn, format string, params ...interface{}) ([]string, error) {
	return call(conn, conn, fmt.Sprintf(format, params...))
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
