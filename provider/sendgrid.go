package provider

import (
	"log"

	"github.com/microapis/messages-api"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// ParamsSendGrid ...
type ParamsSendGrid struct {
	APIKey string `env:"SENDGRID_API_KEY"`
}

// SetParam ...
func (p *ParamsSendGrid) SetParam(key string, value string) {
	switch key {
	case "APIKeys":
		p.APIKey = value
	}
}

// SengridEmailProvider ...
type SengridEmailProvider EmailProvider

// Approve ...
func (p *SengridEmailProvider) Approve(*messages.Email) error {
	return nil
}

// Deliver ...
func (p *SengridEmailProvider) Deliver(m *messages.Email) error {
	// cast params interface
	params := p.Params.(ParamsSendGrid)

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
	client := sendgrid.NewSendClient(params.APIKey)

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
