package email

import (
	"fmt"
	"net/url"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
	// newer version of mailgun
	// mailgun "github.com/mailgun/mailgun-go"
)

const (
	welcomeSubject = "Welcome to gallery.faulkners.io!"
	resetSubject   = "Instructions for resetting your password."
	resetBaseURL   = "https://galleries.faulkners.io/reset"
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

const resetTextTmpl = `Hi there!

It appears that you have requested a password reset. If this was you, please follow the link below to update your password:

%s

If you are asked for a token, please use the following value:

%s

If you didn't request a password reset you can safely ignore this email and your account will not be changed.


Best,

Support Team 
gallery.faulkners.io

`

const resetHTMLTmpl = `Hi there!</br>
</br>
It appears that you have requested a password reset. If this was you, please follow the link 
below to update your password:<br/>
</br>
<a href="%s">%s</a><br/>
</br>
If you are asked for a token, please use the following value:</br>
</br>
%s</br>
</br>
If you didn't request a password reset you can safely ignore 
this email and your account will not be changed.</br>
</br>
Best,</br>
Support Team</br>
gallery.faulkners.io
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

func (c *Client) ResetPw(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message := c.mg.NewMessage(c.from, resetSubject, resetText, toEmail)
	message.SetHtml(resetHTML)

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
