package client

import (
	"fmt"
	"net"
	"time"

	"gopkg.in/fatih/pool.v2"
)

// Client is a client for jcache server
type Client struct {
	connPool pool.Pool
}

// New creates new client instance
func New(addr, user, password string, timeout time.Duration, maxConnections int) (*Client, error) {
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

// Keys returns all keys
func (c *Client) Keys() ([]string, error) {
	conn, err := c.connPool.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	response, err := callConn(conn, "KEYS")
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Get returns value by key
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

	return parseValue(response[0])
}

func callConn(conn net.Conn, format string, params ...interface{}) ([]string, error) {
	return call(conn, conn, fmt.Sprintf(format, params...))
}
