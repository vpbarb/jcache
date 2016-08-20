package client

import (
	"fmt"
	//"log"
	"errors"
	"net"
	"time"

	"gopkg.in/fatih/pool.v2"
)

// Client is a client for jcache server
type Client struct {
	connPool pool.Pool
}

type responseFormatFunc func(response []string) error

var (
	nilResponseFormatter = func(response []string) error { return nil }
)

// New creates new client instance
func New(addr, user, password string, timeout time.Duration, maxConnections int) (*Client, error) {
	factory := func() (net.Conn, error) {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return nil, fmt.Errorf("Cannot connect: %s", err)
		}

		_, err = call(conn, conn, fmt.Sprintf("AUTH %s %s", user, password))
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
func (c *Client) Keys() (keys []string, err error) {
	err = c.call(newKeysResponseFormatter(&keys), 0, "KEYS")
	return
}

// TTL returns ttl of key
func (c *Client) TTL(key string) (ttl time.Duration, err error) {
	err = c.call(newTTLResponseFormatter(&ttl), 1, "TTL %s", key)
	return
}

// Get returns value by key
func (c *Client) Get(key string) (value string, err error) {
	err = c.call(newValueResponseFormatter(&value), 1, "GET %s", key)
	return
}

// Set sets new key value
func (c *Client) Set(key, value string, ttl time.Duration) error {
	return c.call(nilResponseFormatter, 0, `SET %s "%s" %s`, key, value, ttl)
}

// Update updates existing key
func (c *Client) Update(key, value string) error {
	return c.call(nilResponseFormatter, 0, `UPD %s "%s"`, key, value)
}

// Delete deletes value by key
func (c *Client) Delete(key string) error {
	return c.call(nilResponseFormatter, 0, "DEL %s", key)
}

func (c *Client) call(responseFormatter responseFormatFunc, minLines int, command string, params ...interface{}) error {
	conn, err := c.connPool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()

	response, err := call(conn, conn, fmt.Sprintf(command, params...))
	if err != nil {
		return err
	}
	if len(response) < minLines {
		return errors.New("Invalid response rows count")
	}

	return responseFormatter(response)
}

func newKeysResponseFormatter(keys *[]string) responseFormatFunc {
	return func(response []string) (err error) {
		*keys = response
		return
	}
}

func newValueResponseFormatter(value *string) responseFormatFunc {
	return func(response []string) (err error) {
		*value, err = parseValue(response[0])
		return
	}
}

func newTTLResponseFormatter(ttl *time.Duration) responseFormatFunc {
	return func(response []string) (err error) {
		*ttl, err = parseTTL(response[0])
		return
	}
}
