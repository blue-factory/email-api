package provider

import (
	"errors"
	"html"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/microapis/messages-api"
	"github.com/microapis/messages-api/channel"
	"github.com/stoewer/go-strcase"
)

const (
	// SESProvider the provider name
	SESProvider = "ses"
	// SESAWSKeyID the SES key ID
	SESAWSKeyID = "SesAwsKeyId"
	// SESAWSSecretKey the SES secret key
	SESAWSSecretKey = "SesAwsSecretKey"
	// SESAWSRegion the SES secret key
	SESAWSRegion = "SesAwsRegion"
)

// SESEmailProvider ...
type SESEmailProvider channel.Provider

// NewSES ...
func NewSES() *SESEmailProvider {
	p := &SESEmailProvider{
		Name:   SESProvider,
		Params: make(map[string]string),
	}

	p.Params[SESAWSKeyID] = ""
	p.Params[SESAWSSecretKey] = ""
	p.Params[SESAWSRegion] = ""

	return p
}

// Approve ...
func (p *SESEmailProvider) Approve(*messages.Email) error {
	return nil
}

// Deliver ...
func (p *SESEmailProvider) Deliver(m *messages.Email) error {
	// define aws config credentials
	config := &aws.Config{
		Region:      aws.String(p.Params[SESAWSRegion]),
		Credentials: credentials.NewStaticCredentials(p.Params[SESAWSKeyID], p.Params[SESAWSSecretKey], ""),
	}

	// define aws session
	session := session.New(config)

	// define ses instance
	sesClient := ses.New(session)

	// if not has HTML, set text message to HTML
	if m.HTML == "" {
		m.HTML = m.Text
	}

	// unescape html source
	m.HTML = html.UnescapeString(m.HTML)

	// prepare message with email values
	msg := &ses.Message{
		Subject: &ses.Content{
			Charset: aws.String("utf-8"),
			Data:    &m.Subject,
		},
		Body: &ses.Body{
			Html: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &m.HTML,
			},
			Text: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    &m.Text,
			},
		},
	}

	// define emails destinations
	dest := &ses.Destination{
		ToAddresses: aws.StringSlice(m.To),
	}

	// send emails to destinations
	_, err := sesClient.SendEmail(&ses.SendEmailInput{
		Source:           &m.From,
		Destination:      dest,
		Message:          msg,
		ReplyToAddresses: aws.StringSlice(m.ReplyTo),
	})
	if err != nil {
		return err
	}

	return nil
}

// LoadEnv ...
func (p *SESEmailProvider) LoadEnv() error {
	env := strings.ToUpper(strcase.SnakeCase(SESAWSKeyID))
	value := os.Getenv("PROVIDER_" + env)
	if value == "" {
		return errors.New("PROVIDER_" + env + " env value not defined")
	}

	p.Params[SESAWSKeyID] = value

	env = strings.ToUpper(strcase.SnakeCase(SESAWSSecretKey))
	value = os.Getenv("PROVIDER_" + env)
	if value == "" {
		return errors.New("PROVIDER_" + env + " env value not defined")
	}

	p.Params[SESAWSSecretKey] = value

	env = strings.ToUpper(strcase.SnakeCase(SESAWSRegion))
	value = os.Getenv("PROVIDER_" + env)
	if value == "" {
		return errors.New("PROVIDER_" + env + " env value not defined")
	}

	p.Params[SESAWSRegion] = value

	return nil
}

// ToProvider ...
func (p *SESEmailProvider) ToProvider() *channel.Provider {
	return &channel.Provider{
		Name:   p.Name,
		Params: p.Params,
	}
}
