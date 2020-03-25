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

	// get redis host env value
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		err := errors.New("invalid REDIS_HOST env value")
		log.Fatal(err)
	}

	// get redis port env value
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		err := errors.New("invalid REDIS_PORT env value")
		log.Fatal(err)
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// get redis database env value
	redisDatabase := os.Getenv("REDIS_DATABASE")
	if redisDatabase == "" {
		err := errors.New("invalid REDIS_DATABASE env value")
		log.Fatal(err)
	}

	// read providers env values
	providersEnv := os.Getenv("PROVIDERS")
	if providersEnv == "" {
		log.Fatal(errors.New("PROVIDERS value not defined"))
	}
	// define provider slice names
	providers := strings.Split(providersEnv, ",")

	email.Run(addr, redisAddr, redisDatabase, providers)
}
