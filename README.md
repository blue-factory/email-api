# Messages Email API

The microservice that extends from [Messages API](https://github.com/microapis/messages-api) is responsible for sending emails through several providers such as Sendgrid, Mandrill or AWS SES.

As explained in the Messages API repository, it can be seen that there are three models, messages, channel and provider. To know more you can read the readme of Messages API. [[Link]](https://github.com/microapis/messages-hook-api)

## Channels

Corresponds to the type of notification that could be sent, for this there must be the implementation of that channel as a gRPC service through an API.

## Providers

The provider corresponds and attribute of **channel** and allows to identify what types of messages are available for a specific channel.

In this api we will find the implementation of only 3 providers. Find **Providers** implementation at the [`./provider`](./lib) folder.

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
