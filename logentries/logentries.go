package logentries

import (
	"crypto/tls"
	"fmt"
)

const DefaultUrl = "data.logentries.com:443"

type Client struct {
	conn  *tls.Conn
	token string
	url   string
}

func New(url, token string) (*Client, error) {
	c := &Client{token: token, url: url}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Write(b []byte) (int, error) {
	s := fmt.Sprintf("%s %s\n", c.token, b)
	return c.writeAndRetry([]byte(s))
}

func (c *Client) connect() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	conn, err := tls.Dial("tcp", "data.logentries.com:443", &tls.Config{})
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Client) writeAndRetry(b []byte) (int, error) {
	if c.conn != nil {
		if n, err := c.write(b); err == nil {
			return n, err
		}
	}
	if err := c.connect(); err != nil {
		return 0, err
	}
	return c.write(b)
}
