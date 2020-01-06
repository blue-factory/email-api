package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/microapis/messages-api"
	"github.com/microapis/messages-email-api/provider"
)

// Backend ...
type Backend struct {
	Sendgrid *provider.SendgridProvider
	Mandrill *provider.MandrillProvider
	Ses      *provider.SESProvider
}

// NewBackend ...
func NewBackend(sendgrid *provider.SendgridProvider, mandrill *provider.MandrillProvider, ses *provider.SESProvider) *Backend {
	return &Backend{
		Sendgrid: sendgrid,
		Mandrill: mandrill,
		Ses:      ses,
	}
}

// Approve ...
func (s *Backend) Approve(content string) (valid bool, err error) {
	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Request] content = %v", content))

	if content == "" {
		err := errors.New("Invalid message content")
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Error] error = %v", err))
		return false, err
	}

	m := new(messages.Email)
	err = json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Error] error = %v", err))
		return false, err
	}

	switch m.Provider {
	case provider.SendgridName:
		err = s.Sendgrid.Approve(m)
	case provider.MandrillName:
		err = s.Mandrill.Approve(m)
	case provider.SESName:
		err = s.Ses.Approve(m)
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Error] error = %v", err))
		return false, err
	}

	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Approve][Response] message = %v", m))

	return true, nil
}

// Deliver ...
func (s *Backend) Deliver(content string) error {
	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Request] content = %v", content))

	if content == "" {
		err := errors.New("Invalid message content")
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Error] error = %v", err))
		return err
	}

	m := new(messages.Email)
	err := json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Error] error = %v", err))
		return err
	}

	switch m.Provider {
	case provider.SendgridName:
		err = s.Sendgrid.Deliver(m)
	case provider.MandrillName:
		err = s.Mandrill.Deliver(m)
	case provider.SESName:
		err = s.Ses.Deliver(m)
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Error] error = %v", err))
		return err
	}

	log.Println(fmt.Sprintf("[gRPC][MessagesEmailService][Deliver][Response] message = %v", m))

	return nil
}
