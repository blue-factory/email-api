package provider

import (
	"html"

	"github.com/microapis/messages-api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// ParamsSES ...
type ParamsSES struct {
	AwsKeyID     string `env:"AWS_KEY_ID"`
	AwsSecretKey string `env:"AWS_SECRET_KEY"`
	AwsRegion    string `env:"AWS_REGION"`
}

// SetParam ...
func (p *ParamsSES) SetParam(key string, value string) {
	switch key {
	case "AwsKeyID":
		p.AwsKeyID = value
	case "AwsSecretKey":
		p.AwsSecretKey = value
	case "AwsRegion":
		p.AwsRegion = value
	}
}

// SESEmailProvider ...
type SESEmailProvider EmailProvider

// Approve ...
func (p *SESEmailProvider) Approve(*messages.Email) error {
	return nil
}

// Deliver ...
func (p *SESEmailProvider) Deliver(m *messages.Email) error {
	// cast params interface
	params := p.Params.(ParamsSES)

	// define aws config credentials
	config := &aws.Config{
		Region:      aws.String(params.AwsRegion),
		Credentials: credentials.NewStaticCredentials(params.AwsKeyID, params.AwsSecretKey, ""),
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
