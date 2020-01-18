package provider

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/microapis/messages-core/channel"
	messagesemail "github.com/microapis/messages-email-api"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stoewer/go-strcase"
)

const (
	// SendgridName the provider name
	SendgridName = "sendgrid"
	// SendgridAPIKey the sendgrid api key
	SendgridAPIKey = "SendgridApiKey"
)

// SendgridProvider ...
type SendgridProvider struct {
	Root channel.Provider
}

// NewSendgrid ...
func NewSendgrid() *SendgridProvider {
	p := &SendgridProvider{
		Root: channel.Provider{
			Name:   SendgridName,
			Params: make(map[string]string),
		},
	}

	p.Root.Params[SendgridAPIKey] = ""

	return p
}

// Keys ...
func (p *SendgridProvider) Keys() []string {
	k := make([]string, 0)
	k = append(k, SendgridAPIKey)
	return k
}

// Approve ...
func (p *SendgridProvider) Approve(*messagesemail.Message) error {
	return nil
}

// Deliver ...
func (p *SendgridProvider) Deliver(m *messagesemail.Message) error {
	// define from and to values
	from := mail.NewEmail(m.FromName, m.From)
	to := mail.NewEmail(m.To[0], m.To[0])

	// if not has HTML, set text message to HTML
	if m.HTML == "" {
		m.HTML = m.Text
	}

	// prepare single email
	message := mail.NewSingleEmail(from, m.Subject, to, m.Text, m.HTML)

	// create client
	client := sendgrid.NewSendClient(p.Root.Params[SendgridAPIKey])

	// send message
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	log.Println(response.StatusCode)
	log.Println(response.Body)
	log.Println(response.Headers)

	return nil
}

// LoadEnv ...
func (p *SendgridProvider) LoadEnv() error {
	env := strings.ToUpper(strcase.SnakeCase(SendgridAPIKey))
	value := os.Getenv("PROVIDER_" + env)
	if value == "" {
		return errors.New("PROVIDER_" + env + " env value not defined")
	}

	p.Root.Params[SendgridAPIKey] = value

	return nil
}

// ToProvider ...
func (p *SendgridProvider) ToProvider() *channel.Provider {
	return &channel.Provider{
		Name:   p.Root.Name,
		Params: p.Root.Params,
	}
}
