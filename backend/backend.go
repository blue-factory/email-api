package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	messagesemail "github.com/microapis/email-api"
	"github.com/microapis/email-api/provider"
)

var (
	sendgrid *provider.SendgridProvider
	mandrill *provider.MandrillProvider
	ses      *provider.SESProvider
)

// Backend ...
type Backend struct {
	Sendgrid *provider.SendgridProvider
	Mandrill *provider.MandrillProvider
	Ses      *provider.SESProvider
}

// NewBackend ...
func NewBackend(providers []string) (*Backend, error) {
	var err error

	// iterate over providers name
	for _, v := range providers {
		switch v {
		case provider.SESName:
			ses = provider.NewSES()
			err = ses.LoadEnv()
		case provider.SendgridName:
			sendgrid = provider.NewSendgrid()
			err = sendgrid.LoadEnv()
		case provider.MandrillName:
			mandrill = provider.NewMandrill()
			err = mandrill.LoadEnv()
		}
	}
	if err != nil {
		return nil, err
	}

	return &Backend{
		Sendgrid: sendgrid,
		Mandrill: mandrill,
		Ses:      ses,
	}, nil
}

// Approve ...
func (b *Backend) Approve(content string) (valid bool, err error) {
	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Request] content = %v", content))

	if content == "" {
		err := errors.New("Invalid message content")
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Error] error = %v", err))
		return false, err
	}

	m := new(messagesemail.Message)
	err = json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Error] error = %v", err))
		return false, err
	}

	switch m.Provider {
	case provider.SendgridName:
		err = b.Sendgrid.Approve(m)
	case provider.MandrillName:
		err = b.Mandrill.Approve(m)
	case provider.SESName:
		err = b.Ses.Approve(m)
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Error] error = %v", err))
		return false, err
	}

	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Response] message = %v", m))

	return true, nil
}

// Deliver ...
func (b *Backend) Deliver(content string) error {
	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Request] content = %v", content))

	if content == "" {
		err := errors.New("Invalid message content")
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Error] error = %v", err))
		return err
	}

	m := new(messagesemail.Message)
	err := json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Error] error = %v", err))
		return err
	}

	switch m.Provider {
	case provider.SendgridName:
		err = b.Sendgrid.Deliver(m)
	case provider.MandrillName:
		err = b.Mandrill.Deliver(m)
	case provider.SESName:
		err = b.Ses.Deliver(m)
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Error] error = %v", err))
		return err
	}

	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Response] message = %v", m))

	return nil
}
