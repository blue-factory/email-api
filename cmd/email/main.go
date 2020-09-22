package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	email "github.com/microapis/email-api/run"
)

func main() {
	// get grpc host env value
	host := os.Getenv("HOST")
	if host == "" {
		err := errors.New("invalid HOST env value")
		log.Fatal(err)
	}

	// get grpc port env value
	port := os.Getenv("PORT")
	if port == "" {
		err := errors.New("invalid PORT env value")
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	// get redis url env value
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		err := errors.New("invalid REDIS_URL env value")
		log.Fatal(err)
	}

	// read providers env values
	providersEnv := os.Getenv("PROVIDERS")
	if providersEnv == "" {
		log.Fatal(errors.New("PROVIDERS value not defined"))
	}
	// define provider slice names
	providers := strings.Split(providersEnv, ",")

	email.Run(addr, redisURL, providers)
}
