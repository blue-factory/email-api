# Messages Email API

Microservice implemented from [Messages Core](https://github.com/microapis/messages-core) is responsible for sending emails through several providers such as Sendgrid, Mandrill or AWS SES.

As explained in the Messages Core repository, it can be seen that there are three models, messages, channel and provider. To know more you can read the readme of Messages Core.

## Channels

Corresponds to the type of notification that could be sent, for this there must be the implementation of that channel as a gRPC service through an API.

## Providers

The provider is an attribute of **channel** and allows to identify what types of messages are available for a specific channel.

In this api we will find the implementation of only 3 providers. Find **Providers** implementation at the [`./provider`](./providers) folder.

| Name                                   | Description                    | ENV (each with prefix `PROVIDER_`)                               |
| -------------------------------------- | ------------------------------ | ---------------------------------------------------------------- |
| [Sendgrid](https://sendgrid.com/)      | Free Send 40,000 month.        | `SENDGRID_API_KEY`: string                                       |
| [Mandrill](https://mandrill.com/)      | Free Send 10,000 month.        | `MANDRIL_API_KEY`: string                                        |
| [AWS SES](https://aws.amazon.com/ses/) | \$0.10 for every 1,000 emails. | `SES_AWS_KEY_ID`, `SES_AWS_SECRET_KEY`, `SES_AWS_REGION`: string |

## gRPC Service

```go
service MessageBackendService {
  rpc Approve(MessageBackendApproveRequest) returns (MessageBackendApproveResponse) {}
  rpc Deliver(MessageBackendDeliverRequest) returns (MessageBackendDeliverResponse) {}
}
```

## Model

```go
Email {
  from:       string
  from_name:  string
  to:         []string
  reply_to:   []string
  subject:    string
  text:       string
  html:       string
  provider:   string  // sendgrid, mandrill, ses
}
```

## Commands (Development)

`make build`: build user service for osx.

`make linux`: build user service for linux os.

`make docker`: build docker.

`make r`: run service.

`docker run -it -p 5050:5050 messages-email-api`: run docker.

**Run messages service:**

```sh
HOST=<host> \
PORT=<port> \
MESSAGES_HOST=<messages-host> \
MESSAGES_PORT=<messages-port> \
PROVIDERS=<providers> \
PROVIDER_SENDGRID_API_KEY=<> \
PROVIDER_MANDRIL_API_KEY=<> \
PROVIDER_SES_AWS_KEY_ID=<> \
PROVIDER_SES_AWS_SECRET_KEY=<> \
PROVIDER_SES_AWS_REGION=<> \
./bin/messages-api
```

## TODO

- [ ] Task 1.
- [ ] Task 2.
- [ ] Task 3.
