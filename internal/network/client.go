package network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	addr    string
	timeout time.Duration
}

func NewClient(addr string) *Client {
	return &Client{addr: addr, timeout: 5 * time.Second}
}

func (c *Client) SendRequest(rawCmd string) (string, error) {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(c.timeout))

	parts := strings.Fields(rawCmd)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	var respReq strings.Builder
	if _, err := fmt.Fprintf(&respReq, "*%d\r\n", len(parts)); err != nil {
		return "", err
	}

	for _, part := range parts {
		if _, err := fmt.Fprintf(&respReq, "$%d\r\n%s\r\n", len(part), part); err != nil {
			return "", err
		}
	}

	_, err = conn.Write([]byte(respReq.String()))
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)
	return ParseResponse(reader)
}
