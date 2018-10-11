package email

import (
	"fmt"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

func WithMailgun(domain, apiKey, publicApiKey string) ClientConfig {
	return func(c *Client) {
		c.mg = mailgun.NewMailgun(domain, apiKey, publicApiKey)
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// Set a default from email address...
		from: "support@mail.jackytck.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
