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
NAME=email
SVC=$(NAME)-api
BIN_PATH=$(PWD)/bin
BIN=$(BIN_PATH)/$(SVC)
REGISTRY_URL=$(DOCKER_USER)

HOST=localhost
PORT=5030

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DATABASE=1
PROVIDERS=sendgrid,mandrill,ses
PROVIDER_SENDGRID_API_KEY=<PROVIDER_SENDGRID_API_KEY>
PROVIDER_MANDRILL_API_KEY=456
PROVIDER_SES_AWS_KEY_ID=789
PROVIDER_SES_AWS_SECRET_KEY=012
PROVIDER_SES_AWS_REGION=us-east-1

clean c:
	@echo "[clean] Cleaning bin folder..."
	@rm -rf bin/

run r:
	@echo "[running] Running service..."
	@HOST=$(HOST) \
	 PORT=$(PORT) \
	 REDIS_HOST=$(REDIS_HOST) \
	 REDIS_PORT=$(REDIS_PORT) \
	 REDIS_DATABASE=$(REDIS_DATABASE) \
	 PROVIDERS=$(PROVIDERS)	\
	 PROVIDER_SENDGRID_API_KEY=$(PROVIDER_SENDGRID_API_KEY) \
	 PROVIDER_MANDRILL_API_KEY=$(PROVIDER_MANDRILL_API_KEY) \
	 PROVIDER_SES_AWS_KEY_ID=$(PROVIDER_SES_AWS_KEY_ID) \
	 PROVIDER_SES_AWS_SECRET_KEY=$(PROVIDER_SES_AWS_SECRET_KEY) \
	 PROVIDER_SES_AWS_REGION=$(PROVIDER_SES_AWS_REGION) \
	 go run cmd/$(NAME)/main.go

build b: proto
	@echo "[build] Building service..."
	@cd cmd/$(NAME) && go build -o $(BIN)

linux l:
	@echo "[build-linux] Building service..."
	@cd cmd/$(NAME) && GOOS=linux GOARCH=amd64 go build -o $(BIN)

add-migration am: 
	@echo "[add-migration] Adding migration"
	@goose -dir "./database/migrations" create $(name) sql

migrations m:
	@echo "[migrations] Runing migrations..."
	@cd database/migrations && goose postgres $(DSN) up

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

clean-proto cp:
	@echo "[clean-proto] Cleaning proto files..."
	@rm -rf proto/*.pb.go || true

proto pro: clean-proto
	@echo "[proto] Generating proto file..."
	@protoc -I proto -I $(GOPATH)/src --go_out=plugins=grpc:./proto ./proto/*.proto 

test t:
	@echo "[test] Testing $(NAME)..."
	@HOST=$(HOST) \
	 PORT=$(PORT) \
	 POSTGRES_DSN=$(POSTGRES_DSN) \
	 go test -count=1 -v ./client/$(NAME)_test.go

.PHONY: clean c run r build b linux l add-migration am migrations m docker d docker-login dl push p compose co stop s clean-proto cp proto pro test t