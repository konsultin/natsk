package natsk

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type Client struct {
	Conn *nats.Conn
}

type MsgHandler func(msg *nats.Msg)

func New(url string) (*Client, error) {
	// Add some default options like ReconnectWait and MaxReconnects for robustness
	opts := []nats.Option{
		nats.Name("Konsultin API"),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			// We could log here if we had a logger passed in, but for now we rely on the main app logging
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			// Reconnected
		}),
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("natsk: failed to connect: %w", err)
	}

	return &Client{Conn: nc}, nil
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Drain()
		c.Conn.Close()
	}
}

func (c *Client) Subscribe(subject string, handler MsgHandler) (*nats.Subscription, error) {
	sub, err := c.Conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return nil, fmt.Errorf("natsk: subscribe failed: %w", err)
	}
	return sub, nil
}

func (c *Client) QueueSubscribe(subject, queue string, handler MsgHandler) (*nats.Subscription, error) {
	sub, err := c.Conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		handler(msg)
	})
	if err != nil {
		return nil, fmt.Errorf("natsk: queue subscribe failed: %w", err)
	}
	return sub, nil
}

func (c *Client) Publish(subject string, data []byte) error {
	if err := c.Conn.Publish(subject, data); err != nil {
		return fmt.Errorf("natsk: publish failed: %w", err)
	}
	return nil
}

func IsConnectionError(err error) bool {
	// Check if error is related to connection issues
	// specific logic can be added here if needed
	return strings.Contains(err.Error(), "no servers available") ||
		err == nats.ErrNoServers
}
