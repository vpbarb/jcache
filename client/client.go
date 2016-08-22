package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Barberrrry/jcache/protocol"
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
func (c *Client) Keys() ([]string, error) {
	request := protocol.NewKeysRequest()
	response := protocol.NewKeysResponse()
	if err := c.call(request, response); err != nil {
		return nil, err
	}

	return response.Keys, response.Error
}

// Get returns value by key
func (c *Client) Get(key string) (string, error) {
	request := protocol.NewGetRequest()
	request.Key = key
	response := protocol.NewGetResponse()
	if err := c.call(request, response); err != nil {
		return "", err
	}

	return response.Value, response.Error
}

// Set sets new key value
func (c *Client) Set(key, value string, ttl uint64) error {
	request := protocol.NewSetRequest()
	request.Key = key
	request.Value = value
	request.TTL = ttl
	response := protocol.NewSetResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// Update updates existing key
func (c *Client) Update(key, value string) error {
	request := protocol.NewUpdRequest()
	request.Key = key
	request.Value = value
	response := protocol.NewUpdResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// Delete deletes value by key
func (c *Client) Delete(key string) error {
	request := protocol.NewDelRequest()
	request.Key = key
	response := protocol.NewDelResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// HashCreate creates new hash with ttl
func (c *Client) HashCreate(key string, ttl uint64) error {
	request := protocol.NewHashCreateRequest()
	request.Key = key
	request.TTL = ttl
	response := protocol.NewHashCreateResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// HashSet add new field to hash
func (c *Client) HashSet(key, field, value string) error {
	request := protocol.NewHashSetRequest()
	request.Key = key
	request.Field = field
	request.Value = value
	response := protocol.NewHashSetResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// HashGet returns hash field value
func (c *Client) HashGet(key, field string) (string, error) {
	request := protocol.NewHashGetRequest()
	request.Key = key
	request.Field = field
	response := protocol.NewHashGetResponse()
	if err := c.call(request, response); err != nil {
		return "", err
	}

	return response.Value, response.Error
}

// HashGetAll returns all hash fields values
func (c *Client) HashGetAll(key string) (map[string]string, error) {
	request := protocol.NewHashGetAllRequest()
	request.Key = key
	response := protocol.NewHashGetAllResponse()
	if err := c.call(request, response); err != nil {
		return nil, err
	}

	return response.Fields, response.Error
}

// HashDelete deletes field from hash
func (c *Client) HashDelete(key, field string) error {
	request := protocol.NewHashDelRequest()
	request.Key = key
	request.Field = field
	response := protocol.NewHashDelResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// HashKeys returns all hash fields
func (c *Client) HashKeys(key string) ([]string, error) {
	request := protocol.NewHashKeysRequest()
	request.Key = key
	response := protocol.NewHashKeysResponse()
	if err := c.call(request, response); err != nil {
		return nil, err
	}

	return response.Keys, response.Error
}

// HashLength returns count of hash elements
func (c *Client) HashLength(key string) (int, error) {
	request := protocol.NewHashLenRequest()
	request.Key = key
	response := protocol.NewHashLenResponse()
	if err := c.call(request, response); err != nil {
		return 0, err
	}

	return response.Len, response.Error
}

// ListCreate creates new list with ttl
func (c *Client) ListCreate(key string, ttl uint64) error {
	request := protocol.NewListCreateRequest()
	request.Key = key
	request.TTL = ttl
	response := protocol.NewListCreateResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// ListRightPush adds new value to the list ending
func (c *Client) ListRightPush(key, value string) error {
	request := protocol.NewListRightPushRequest()
	request.Key = key
	request.Value = value
	response := protocol.NewListRightPushResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// ListLeftPush adds new value to the list beginning
func (c *Client) ListLeftPush(key, value string) error {
	request := protocol.NewListLeftPushRequest()
	request.Key = key
	request.Value = value
	response := protocol.NewListLeftPushResponse()
	if err := c.call(request, response); err != nil {
		return err
	}

	return response.Error
}

// ListRightPop returns and removes the value from the list ending
func (c *Client) ListRightPop(key string) (string, error) {
	request := protocol.NewListRightPopRequest()
	request.Key = key
	response := protocol.NewListRightPopResponse()
	if err := c.call(request, response); err != nil {
		return "", err
	}

	return response.Value, response.Error
}

// ListLeftPop returns and removes the value from the list beginning
func (c *Client) ListLeftPop(key string) (string, error) {
	request := protocol.NewListLeftPopRequest()
	request.Key = key
	response := protocol.NewListLeftPopResponse()
	if err := c.call(request, response); err != nil {
		return "", err
	}

	return response.Value, response.Error
}

// ListLength returns count of list elements
func (c *Client) ListLength(key string) (int, error) {
	request := protocol.NewListLenRequest()
	request.Key = key
	response := protocol.NewListLenResponse()
	if err := c.call(request, response); err != nil {
		return 0, err
	}

	return response.Len, response.Error
}

// ListRange returns all list values from start to stop
func (c *Client) ListRange(key string, start, stop int) ([]string, error) {
	request := protocol.NewListRangeRequest()
	request.Key = key
	request.Start = start
	request.Stop = stop
	response := protocol.NewListRangeResponse()
	if err := c.call(request, response); err != nil {
		return nil, err
	}

	return response.Values, response.Error
}

func (c *Client) connFactory() (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect: %s", err)
	}

	request := protocol.NewAuthRequest()
	request.User = c.user
	request.Password = c.password
	response := protocol.NewAuthResponse()

	err = c.callRW(conn, request, response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		conn.Close()
		return nil, fmt.Errorf("Cannot authentiticate: %s", response.Error)
	}
	return conn, nil
}

func (c *Client) call(request protocol.Encoder, response protocol.Decoder) error {
	conn, err := c.connPool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()

	return c.callRW(conn, request, response)
}

func (c *Client) callRW(rw io.ReadWriter, request protocol.Encoder, response protocol.Decoder) error {
	data, err := request.Encode()
	if err != nil {
		return err
	}
	_, err = rw.Write(data)
	if err != nil {
		return err
	}

	rb := bufio.NewReader(rw)
	line, _, err := rb.ReadLine()
	if err != nil {
		return fmt.Errorf("Cannot read from connection: %s", err)
	}
	err = response.Decode(line, rb)
	if err != nil {
		return err
	}
	return nil
}
