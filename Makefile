#
# SO variables
#
# DOCKER_USER
# DOCKER_PASS
#

#
# Internal variables
#
VERSION=0.2.3
SVC=messages-email-api
BIN_PATH=$(PWD)/bin
BIN=$(BIN_PATH)/$(SVC)
REGISTRY_URL=$(DOCKER_USER)

#
# SVC variables
#
PORT=5050
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DATABASE=1
PROVIDERS=sendgrid,mandrill,ses
PROVIDER_SENDGRID_API_KEY=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
PROVIDER_MANDRILL_API_KEY=456
PROVIDER_SES_AWS_KEY_ID=789
PROVIDER_SES_AWS_SECRET_KEY=012
PROVIDER_SES_AWS_REGION=us-east-1

clean c:
	@echo "[clean] Cleaning bin folder..."
	@rm -rf bin/

run r:
	@echo "[running] Running service..."
	@PROVIDERS=$(PROVIDERS) PROVIDER_SENDGRID_API_KEY=$(PROVIDER_SENDGRID_API_KEY) PORT=$(PORT) REDIS_HOST=$(REDIS_HOST) REDIS_PORT=$(REDIS_PORT) REDIS_DATABASE=$(REDIS_DATABASE) PROVIDER_MANDRILL_API_KEY=$(PROVIDER_MANDRILL_API_KEY) PROVIDER_SES_AWS_KEY_ID=$(PROVIDER_SES_AWS_KEY_ID) PROVIDER_SES_AWS_SECRET_KEY=$(PROVIDER_SES_AWS_SECRET_KEY) PROVIDER_SES_AWS_REGION=$(PROVIDER_SES_AWS_REGION) go run cmd/main.go

build b:
	@echo "[build] Building service..."
	@cd cmd && go build -o $(BIN)

linux l: 
	@echo "[build-linux] Building service..."
	@cd cmd && GOOS=linux GOARCH=amd64 go build -o $(BIN)

docker d:
	@echo "[docker] Building image..."
	@docker build -t $(SVC):$(VERSION) .
	
docker-login dl:
	@echo "[docker] Login to docker..."
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)

push p: linux docker docker-login
	@echo "[docker] pushing $(REGISTRY_URL)/$(SVC):$(VERSION)"
	@docker tag $(SVC):$(VERSION) $(REGISTRY_URL)/$(SVC):$(VERSION)
	@docker push $(REGISTRY_URL)/$(SVC):$(VERSION)

compose co:
	@echo "[docker-compose] Running docker-compose..."
	@docker-compose build
	@docker-compose up

stop s: 
	@echo "[docker-compose] Stopping docker-compose..."
	@docker-compose down

.PHONY: clean c run r build b linux l docker d docker-login dl push p compose co stop s