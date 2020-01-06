package provider

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/keighl/mandrill"
	"github.com/microapis/messages-api"
	"github.com/microapis/messages-api/channel"
	"github.com/stoewer/go-strcase"
)

const (
	// MandrillProvider the provider name
	MandrillProvider = "mandrill"
	// MandrillAPIKey the mandrill api key
	MandrillAPIKey = "MandrillApiKey"
)

// MandrillEmailProvider ...
type MandrillEmailProvider channel.Provider

// NewMandrill ...
func NewMandrill() *MandrillEmailProvider {
	p := &MandrillEmailProvider{
		Name:   MandrillProvider,
		Params: make(map[string]string),
	}

	p.Params[MandrillAPIKey] = ""

	return p
}

// Keys ...
func (p *MandrillEmailProvider) Keys() []string {
	k := make([]string, 0)
	k = append(k, MandrillAPIKey)
	return k
}

// Approve ...
func (p *MandrillEmailProvider) Approve(*messages.Email) error {
	return nil
}

// Deliver ...
func (p *MandrillEmailProvider) Deliver(m *messages.Email) error {
	// create client
	client := mandrill.ClientWithKey(p.Params[MandrillAPIKey])

	// if not has HTML, set text message to HTML
	if m.HTML == "" {
		m.HTML = m.Text
	}

	// prepare message
	email := &mandrill.Message{
		FromEmail: m.From,
		FromName:  m.FromName,
		Subject:   m.Subject,
		HTML:      m.HTML,
		Text:      m.Text,
	}

	// add email recipient
	email.AddRecipient(m.To[0], m.To[0], "to")

	// send email
	response, err := client.MessagesSend(email)
	if err != nil {
		return err
	}

	log.Println(response[0].Email)
	log.Println(response[0].Id)
	log.Println(response[0].Status)
	log.Println(response[0].RejectionReason)

	return nil
}

// LoadEnv ...
func (p *MandrillEmailProvider) LoadEnv() error {
	env := strings.ToUpper(strcase.SnakeCase(MandrillAPIKey))
	value := os.Getenv("PROVIDER_" + env)
	if value == "" {
		return errors.New("PROVIDER_" + env + " env value not defined")
	}

	p.Params[MandrillAPIKey] = value

	return nil
}

// ToProvider ...
func (p *MandrillEmailProvider) ToProvider() *channel.Provider {
	return &channel.Provider{
		Name:   p.Name,
		Params: p.Params,
	}
}
