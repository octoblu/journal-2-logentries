package logentries

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"strings"
	"fmt"
	"github.com/octoblu/journal-2-logentries/logline"
)

const DefaultUrl = "api.logentries.com:20000"

type Client struct {
	conn  *tls.Conn
	pool  *x509.CertPool
	token string
	url   string
}

func New(url, token string) (*Client, error) {
	c := &Client{token: token, url: url}
	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(pemCerts); !ok {
		return nil, errors.New("failed to parse certs")
	}
	c.pool = pool
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) Write(logLine logline.LogLine) (int, error) {
	prettyOutput := fmt.Sprintf("%s %s %s", logLine.FormatTimestamp(), logLine.Unit, logLine.Message)
	prettyOutput = strings.Replace(prettyOutput, "\n", "\\n", -1)
	s := fmt.Sprintf("%s %s\n", c.token, prettyOutput)
	return c.writeAndRetry([]byte(s))
}

func (c *Client) connect() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	conn, err := tls.Dial("tcp", c.url, &tls.Config{RootCAs: c.pool})
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
