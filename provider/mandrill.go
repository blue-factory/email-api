package provider

import (
	"log"

	"github.com/microapis/messages-api"

	"github.com/keighl/mandrill"
)

// ParamsMandrill ...
type ParamsMandrill struct {
	APIKey string `env:"MANDRILL_API_KEY"`
}

// SetParam ...
func (p *ParamsMandrill) SetParam(key string, value string) {
	switch key {
	case "APIKeys":
		p.APIKey = value
	}
}

// MandrillEmailProvider ...
type MandrillEmailProvider EmailProvider

// Approve ...
func (p *MandrillEmailProvider) Approve(*messages.Email) error {
	return nil
}

// Deliver ...
func (p *MandrillEmailProvider) Deliver(m *messages.Email) error {
	// cast params interface
	params := p.Params.(ParamsMandrill)

	// create client
	client := mandrill.ClientWithKey(params.APIKey)

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
