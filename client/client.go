package client

import (
	"bufio"
	"fmt"
	"net"

	"gopkg.in/fatih/pool.v2"
)

type Client struct {
	pool pool.Pool
}

func NewClient(addr, user, password string, maxConnections int) (*Client, error) {
	factory := func() (net.Conn, error) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, fmt.Errorf("Cannot connect: %s", err)
		}

		response, err := call(conn, fmt.Sprintf("AUTH %s %s", user, password))
		if err != nil {
			return nil, err
		}

		switch response {
		case "OK":
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

func call(conn net.Conn, command string) (string, error) {
	w := bufio.NewWriter(conn)
	w.WriteString(command + "\n")
	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("Cannot write to connection: %s", err)
	}

	r := bufio.NewReader(conn)
	line, _, err := r.ReadLine()
	if err != nil {
		return "", fmt.Errorf("Cannot read from connection: %s", err)
	}

	return string(line), nil
}

func (c *Client) Get(key string) (string, error) {
	conn, err := c.pool.Get()
	if err != nil {
		return "", err
	}

	response, err := call(conn, fmt.Sprintf("GET %s", key))
	if err != nil {
		return "", err
	}

	return response, nil
}
