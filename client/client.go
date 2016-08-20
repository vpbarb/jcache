package client

import (
	"errors"
	"fmt"
	"net"
	"strconv"
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
	err = c.call(newKeysResponseFormatter(&keys), "KEYS")
	return
}

// TTL returns ttl of key
func (c *Client) TTL(key string) (ttl time.Duration, err error) {
	err = c.call(newTTLResponseFormatter(&ttl), "TTL %s", key)
	return
}

// Get returns value by key
func (c *Client) Get(key string) (value string, err error) {
	err = c.call(newValueResponseFormatter(&value), "GET %s", key)
	return
}

// Set sets new key value
func (c *Client) Set(key, value string, ttl time.Duration) error {
	return c.call(nilResponseFormatter, `SET %s "%s" %s`, key, value, ttl)
}

// Update updates existing key
func (c *Client) Update(key, value string) error {
	return c.call(nilResponseFormatter, `UPD %s "%s"`, key, value)
}

// Delete deletes value by key
func (c *Client) Delete(key string) error {
	return c.call(nilResponseFormatter, "DEL %s", key)
}

// HashCreate creates new hash with ttl
func (c *Client) HashCreate(key string, ttl time.Duration) error {
	return c.call(nilResponseFormatter, "HCREATE %s %s", key, ttl)
}

// HashSet add new field to hash
func (c *Client) HashSet(key, field, value string) error {
	return c.call(nilResponseFormatter, `HSET %s %s "%s"`, key, field, value)
}

// HashGet returns hash field value
func (c *Client) HashGet(key, field string) (value string, err error) {
	err = c.call(newValueResponseFormatter(&value), "HGET %s %s", key, field)
	return
}

// HashGetAll returns all hash fields values
func (c *Client) HashGetAll(key string) (hash map[string]string, err error) {
	hash = make(map[string]string)
	err = c.call(newHashResponseFormatter(hash), "HGETALL %s", key)
	return
}

// HashDelete deletes field from hash
func (c *Client) HashDelete(key, field string) error {
	return c.call(nilResponseFormatter, `HDEL %s %s`, key, field)
}

// HashKeys returns all hash fields
func (c *Client) HashKeys(key string) (keys []string, err error) {
	err = c.call(newKeysResponseFormatter(&keys), "HKEYS %s", key)
	return
}

// HashLength returns count of hash elements
func (c *Client) HashLength(key string) (len int, err error) {
	err = c.call(newLenResponseFormatter(&len), "HLEN %s", key)
	return
}

// ListCreate creates new list with ttl
func (c *Client) ListCreate(key string, ttl time.Duration) error {
	return c.call(nilResponseFormatter, "LCREATE %s %s", key, ttl)
}

// ListRightPush adds new value to the list ending
func (c *Client) ListRightPush(key, value string) error {
	return c.call(nilResponseFormatter, `LRPUSH %s "%s"`, key, value)
}

// ListLeftPush adds new value to the list beginning
func (c *Client) ListLeftPush(key, value string) error {
	return c.call(nilResponseFormatter, `LLPUSH %s "%s"`, key, value)
}

// ListRightPop returns and removes the value from the list ending
func (c *Client) ListRightPop(key string) (value string, err error) {
	err = c.call(newValueResponseFormatter(&value), "LRPOP %s", key)
	return
}

// ListLeftPop returns and removes the value from the list beginning
func (c *Client) ListLeftPop(key string) (value string, err error) {
	err = c.call(newValueResponseFormatter(&value), "LLPOP %s", key)
	return
}

// ListLength returns count of list elements
func (c *Client) ListLength(key string) (len int, err error) {
	err = c.call(newLenResponseFormatter(&len), "LLEN %s", key)
	return
}

// ListRange returns all list values from start to stop
func (c *Client) ListRange(key string, start, stop int) (values []string, err error) {
	err = c.call(newValuesResponseFormatter(&values), "LRANGE %s %d %d", key, start, stop)
	return
}

func (c *Client) call(responseFormatter responseFormatFunc, command string, params ...interface{}) error {
	conn, err := c.connPool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()

	response, err := call(conn, conn, fmt.Sprintf(command, params...))
	if err != nil {
		return err
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
		if len(response) < 1 {
			return errors.New("Invalid response rows count")
		}
		*value, err = parseValue(response[0])
		return
	}
}

func newValuesResponseFormatter(values *[]string) responseFormatFunc {
	return func(response []string) (err error) {
		for _, line := range response {
			v, err := parseValue(line)
			if err != nil {
				return err
			}
			*values = append(*values, v)
		}
		return
	}
}

func newTTLResponseFormatter(ttl *time.Duration) responseFormatFunc {
	return func(response []string) (err error) {
		if len(response) < 1 {
			return errors.New("Invalid response rows count")
		}
		*ttl, err = parseTTL(response[0])
		return
	}
}

func newHashResponseFormatter(hash map[string]string) responseFormatFunc {
	return func(response []string) (err error) {
		for _, line := range response {
			field, value, err := parseHashField(line)
			if err != nil {
				return err
			}
			hash[field] = value
		}
		return
	}
}

func newLenResponseFormatter(length *int) responseFormatFunc {
	return func(response []string) (err error) {
		if len(response) < 1 {
			return errors.New("Invalid response rows count")
		}
		*length, err = strconv.Atoi(response[0])
		return
	}
}
