package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	mc "github.com/microapis/clients-go/messages"
	"github.com/microapis/messages-api"
	"github.com/microapis/messages-api/backend"
	"github.com/microapis/messages-email-api/provider"
)

var (
	ses      *provider.SESEmailProvider
	sendgrid *provider.SengridEmailProvider
	mandrill *provider.MandrillEmailProvider
)

type service struct{}

func (s *service) Approve(content string) (valid bool, err error) {
	if content == "" {
		return false, errors.New("Invalid message content")
	}

	fmt.Println(content)
	m := new(messages.Email)

	err = json.Unmarshal([]byte(content), m)
	if err != nil {
		return false, err
	}

	switch m.Provider {
	case "ses":
		err = ses.Approve(m)
	case "sendgrid":
		err = sendgrid.Approve(m)
	case "mandrill":
		err = mandrill.Approve(m)
	}
	if err != nil {
		return false, err
	}

	// TODO(ca): validate content to avoid breakout
	fmt.Println(m)

	return true, nil
}

func (s *service) Deliver(content string) error {
	if content == "" {
		return errors.New("Invalid message content")
	}

	fmt.Println(content)
	m := new(messages.Email)

	err := json.Unmarshal([]byte(content), m)
	if err != nil {
		return err
	}

	switch m.Provider {
	case "ses":
		err = ses.Deliver(m)
	case "sendgrid":
		err = sendgrid.Deliver(m)
	case "mandrill":
		err = mandrill.Deliver(m)
	}
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var err error

	// read providers env values
	providersEnv := os.Getenv("PROVIDERS")
	if providersEnv == "" {
		log.Fatal(errors.New("PROVIDERS value not defined"))
	}

	// define provider slice names
	ppn := strings.Split(providersEnv, ",")
	// define providers slice
	pp := make([]*messages.Provider, 0)

	// iterate over providers name
	for _, v := range ppn {
		switch v {
		case "ses":
			ses = &provider.SESEmailProvider{
				Name:   v,
				Params: provider.ParamsSES{},
			}
			pp = append(pp, &messages.Provider{Name: v, Params: provider.ParamsSES{}})
			err = provider.Load(ses)
		case "sendgrid":
			sendgrid = &provider.SengridEmailProvider{
				Name:   v,
				Params: provider.ParamsSendGrid{},
			}
			pp = append(pp, &messages.Provider{Name: v, Params: provider.ParamsSendGrid{}})
			err = provider.Load(sendgrid)
		case "mandrill":
			mandrill = &provider.MandrillEmailProvider{
				Name:   v,
				Params: provider.ParamsMandrill{},
			}
			pp = append(pp, &messages.Provider{Name: v, Params: provider.ParamsMandrill{}})
			err = provider.Load(mandrill)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	// read messages-api env values
	messagesHost := os.Getenv("MESSAGES_HOST")
	if messagesHost == "" {
		log.Fatal(errors.New("MESSAGES_HOST value not defined"))
	}
	messagesPort := os.Getenv("MESSAGES_PORT")
	if messagesPort == "" {
		log.Fatal(errors.New("MESSAGES_PORT value not defined"))
	}

	// create channel to register
	c := &messages.Channel{
		Name:      "email",
		Host:      messagesHost,
		Port:      messagesPort,
		Providers: pp,
	}

	// register channel on messages-api
	addr := fmt.Sprintf("%s:%s", messagesHost, messagesPort)
	MessagesAPI := mc.NewService(addr)
	err = MessagesAPI.Register(c)
	if err != nil {
		log.Fatalln(err)
	}

	// async checker connection
	go func() {

	}()

	// get grpc port env value:
	port := os.Getenv("PORT")
	if port == "" {
		err := errors.New("invalid PORT env value")
		log.Fatal(err)
	}

	// define address value to grpc service
	addr = fmt.Sprintf("0.0.0.0:%s", port)

	// define service with Approve and Deliver methods
	svc := &service{}

	// start grpc pigeon-ses-api service
	log.Printf("Serving at %s", addr)
	if err := backend.ListenAndServe(addr, svc); err != nil {
		log.Fatal(err)
	}
}
