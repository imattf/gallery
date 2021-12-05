package email

import (
	"fmt"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
	// newer version of mailgun
	// mailgun "github.com/mailgun/mailgun-go"
)

const (
	welcomeSubject = "Welcome to gallery.faulkners.io!"
)

const welcomeText = `Hi there!

Welcome to gallery.faulkners.io! We really hope you enjoy using
our application!

Best,
Matthew
`

const welcomeHTML = `Hi there!<br/>
<br/>
Welcome to
<a href="https://gallery.faulkners.io"> gallery.faulkners.io </a>! We really hope
you enjoy our application!<br/>
<br/>
Best,<br/>
Matthew
`

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

// // works with newer version of mailgun
// func WithMailgun(domain, apiKey string) ClientConfig {
// 	return func(c *Client) {
// 		mg := mailgun.NewMailgun(domain, apiKey)
// 		c.mg = mg
// 	}
// }

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// Set a default email address
		from: "support@faulkners.io",
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
	message := c.mg.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))

	message.SetHtml(welcomeHTML)

	// //works with new version of mailgun...
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	// defer cancel()
	// _, _, err := c.mg.Send(ctx, message)

	_, _, err := c.mg.Send(message)
	if err != nil {
		fmt.Println("Got a mailgun Email error!!", err)
	}
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
