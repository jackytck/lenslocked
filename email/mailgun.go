package email

import (
	"fmt"
	"net/url"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to lenslocked.jackytck.com"
	resetSubject   = "Instructions for resetting your password"
	resetBaseURL   = "https://lenslocked.jackytck.com/reset"
)

const welcomeText = `Hi there!

Welcome to LensLocked! We really hope you enjoy using
our application!

Best,
Jacky
`

const resetTextTmpl = `Hi there!

It appears that you have requested a password reset. If this was you, please follow the link below to update your password:

%s

If you are asked for a token, please use the following value:

%s

If you didn't request a password reset you can safely ignore this email and your account will not be changed.

Best,
LensLocked Support
`

const resetHTMLTmpl = `Hi there!<br/>
<br/>
It appears that you have requested a password reset. If this was you, please follow the link below to update your password:<br/>
<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token, please use the following value:<br/>
<br/>
%s<br/>
<br/>
If you didn't request a password reset you can safely ignore this email and your account will not be changed.<br/>
<br/>
Best,<br/>
LensLocked Support<br/>
`

const welcomeHTML = `Hi there!<br/>
<br/>
Welcome to
<a href="https://lenslocked.jackytck.com">LensLocked</a>! We really hope you enjoy using our application!</br>
<br/>
Best,<br/>
Jacky
`

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

func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()

	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, token, token)

	message := mailgun.NewMessage(c.from, resetSubject, resetText, toEmail)
	message.SetHtml(resetHTML)
	_, _, err := c.mg.Send(message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
