package client

import (
	"fmt"
	"net"
	"time"

	"gopkg.in/fatih/pool.v2"
)

// Client is a client for jcache server
type Client struct {
	addr     string
	timeout  time.Duration
	user     string
	password string
	connPool pool.Pool
}

// New creates new client instance
func New(addr, user, password string, timeout time.Duration, maxConnections int) (*Client, error) {
	client := &Client{
		addr:     addr,
		user:     user,
		password: password,
		timeout:  timeout,
	}

	if connPool, err := pool.NewChannelPool(0, maxConnections, client.connFactory); err == nil {
		client.connPool = connPool
		return client, nil
	} else {
		return nil, fmt.Errorf("Cannot create connection pool: %s", err)
	}
}

// Keys returns all keys
func (c *Client) Keys() (keys []string, err error) {
	err = c.call(keysDataFormatter(&keys), "KEYS\r\n")
	return
}

// TTL returns ttl of key
func (c *Client) TTL(key string) (ttl uint64, err error) {
	err = c.call(ttlDataFormatter(&ttl), "TTL %s\r\n", key)
	return
}

// Get returns value by key
func (c *Client) Get(key string) (value string, err error) {
	err = c.call(valueDataFormatter(&value), "GET %s\r\n", key)
	return
}

// Set sets new key value
func (c *Client) Set(key, value string, ttl uint64) error {
	return c.call(emptyDataFormatter(), "SET %s %d %d\r\n%s\r\n", key, ttl, len(value), value)
}

// Update updates existing key
func (c *Client) Update(key, value string) error {
	return c.call(emptyDataFormatter(), "UPD %s %d\r\n%s\r\n", key, len(value), value)
}

// Delete deletes value by key
func (c *Client) Delete(key string) error {
	return c.call(emptyDataFormatter(), "DEL %s\r\n", key)
}

// HashCreate creates new hash with ttl
func (c *Client) HashCreate(key string, ttl uint64) error {
	return c.call(emptyDataFormatter(), "HCREATE %s %d\r\n", key, ttl)
}

// HashSet add new field to hash
func (c *Client) HashSet(key, field, value string) error {
	return c.call(emptyDataFormatter(), "HSET %s %s %d\r\n%s\r\n", key, field, len(value), value)
}

// HashGet returns hash field value
func (c *Client) HashGet(key, field string) (value string, err error) {
	err = c.call(valueDataFormatter(&value), "HGET %s %s\r\n", key, field)
	return
}

// HashGetAll returns all hash fields values
func (c *Client) HashGetAll(key string) (hash map[string]string, err error) {
	hash = make(map[string]string)
	err = c.call(hashDataFormatter(hash), "HGETALL %s\r\n", key)
	return
}

// HashDelete deletes field from hash
func (c *Client) HashDelete(key, field string) error {
	return c.call(emptyDataFormatter(), "HDEL %s %s\r\n", key, field)
}

// HashKeys returns all hash fields
func (c *Client) HashKeys(key string) (keys []string, err error) {
	err = c.call(keysDataFormatter(&keys), "HKEYS %s\r\n", key)
	return
}

// HashLength returns count of hash elements
func (c *Client) HashLength(key string) (len uint64, err error) {
	err = c.call(lenDataFormatter(&len), "HLEN %s\r\n", key)
	return
}

// ListCreate creates new list with ttl
func (c *Client) ListCreate(key string, ttl uint64) error {
	return c.call(emptyDataFormatter(), "LCREATE %s %d\r\n", key, ttl)
}

// ListRightPush adds new value to the list ending
func (c *Client) ListRightPush(key, value string) error {
	return c.call(emptyDataFormatter(), "LRPUSH %s %d\r\n%s\r\n", key, len(value), value)
}

// ListLeftPush adds new value to the list beginning
func (c *Client) ListLeftPush(key, value string) error {
	return c.call(emptyDataFormatter(), "LLPUSH %s %d\r\n%s\r\n", key, len(value), value)
}

// ListRightPop returns and removes the value from the list ending
func (c *Client) ListRightPop(key string) (value string, err error) {
	err = c.call(valueDataFormatter(&value), "LRPOP %s\r\n", key)
	return
}

// ListLeftPop returns and removes the value from the list beginning
func (c *Client) ListLeftPop(key string) (value string, err error) {
	err = c.call(valueDataFormatter(&value), "LLPOP %s\r\n", key)
	return
}

// ListLength returns count of list elements
func (c *Client) ListLength(key string) (len uint64, err error) {
	err = c.call(lenDataFormatter(&len), "LLEN %s\r\n", key)
	return
}

// ListRange returns all list values from start to stop
func (c *Client) ListRange(key string, start, stop int) (values []string, err error) {
	err = c.call(valuesDataFormatter(&values), "LRANGE %s %d %d\r\n", key, start, stop)
	return
}

func (c *Client) connFactory() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect: %s", err)
	}

	err = transfer(conn, conn, fmt.Sprintf("AUTH %s %s\r\n", c.user, c.password), emptyDataFormatter())
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Cannot authentiticate: %s", err)
	}

	return conn, nil
}

func (c *Client) call(dataFormatter dataFormatFunc, command string, params ...interface{}) error {
	conn, err := c.connPool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()

	return transfer(conn, conn, fmt.Sprintf(command, params...), dataFormatter)
}
