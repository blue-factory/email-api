package run

import (
	"log"

	messagesemail "github.com/microapis/email-api"
	backendEmail "github.com/microapis/email-api/backend"
	"github.com/microapis/messages-core/service"
)

// Run ...
func Run(address string, redisURL string, providers []string) {
	// initialice message backend with Approve and Deliver methods
	backend, err := backendEmail.NewBackend(providers)
	if err != nil {
		log.Fatal(err)
	}

	// initialize message service
	svc, err := service.NewMessageService(messagesemail.Channel, service.ServiceConfig{
		Addr: address,

		RedisURL: redisURL,

		Approve: backend.Approve,
		Deliver: backend.Deliver,
	})
	if err != nil {
		log.Fatal(err)
	}

	// run message service
	err = svc.Run()
	if err != nil {
		log.Fatal(err)
	}
}
