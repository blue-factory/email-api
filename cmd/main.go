package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	mc "github.com/microapis/clients-go/messages"
	"github.com/microapis/messages-api/backend"
	"github.com/microapis/messages-api/channel"
	backendEmail "github.com/microapis/messages-email-api/backend"
	"github.com/microapis/messages-email-api/provider"
)

var (
	ses      *provider.SESProvider
	sendgrid *provider.SendgridProvider
	mandrill *provider.MandrillProvider
)

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
	pp := make([]*channel.Provider, 0)

	// iterate over providers name
	for _, v := range ppn {
		switch v {
		case provider.SESName:
			ses = provider.NewSES()
			err = ses.LoadEnv()
			pp = append(pp, ses.ToProvider())
		case provider.SendgridName:
			sendgrid = provider.NewSendgrid()
			err = sendgrid.LoadEnv()
			pp = append(pp, sendgrid.ToProvider())
		case provider.MandrillName:
			mandrill = provider.NewMandrill()
			err = mandrill.LoadEnv()
			pp = append(pp, mandrill.ToProvider())
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

	// register channel on messages-api
	addr := fmt.Sprintf("%s:%s", messagesHost, messagesPort)

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal(errors.New("HOST value not defined"))
	}

	// get grpc port env value:
	port := os.Getenv("PORT")
	if port == "" {
		err := errors.New("invalid PORT env value")
		log.Fatal(err)
	}

	// create channel to register
	c := &channel.Channel{
		Name:      "email",
		Host:      host,
		Port:      port,
		Providers: pp,
	}

	MessagesAPI := mc.NewService(addr)
	err = MessagesAPI.Register(c)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Channel email is registered on messages-api with providers:", c.ProvidersNames())

	// define address value to grpc service
	addr = fmt.Sprintf("0.0.0.0:%s", port)

	// define service with Approve and Deliver methods
	svc := backendEmail.NewBackend(sendgrid, mandrill, ses)

	// start grpc pigeon-ses-api service
	log.Printf("Serving at %s", addr)
	if err := backend.ListenAndServe(addr, svc); err != nil {
		log.Fatal(err)
	}
}
