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

## Service methods

```go
func (c *Client) Send(e *email.Message, delay int64) (string, error)
func (c *Client) Get(ID string) (*email.Message, error)
func (c *Client) Update(ID string, e *email.Message) error
func (c *Client) Cancel(ID string) error
```

## Environments Values

`PORT`: define email service port.

`HOST`: define email service host.

`REDIS_HOST`: define redis host.

`REDIS_PORT`: define redis port.

`REDIS_DATABASE`: define redis database number.

`PROVIDERS`: define a []string of provider's names

`PROVIDER_SENDGRID_API_KEY`: define sendgrid provider api key.

`PROVIDER_MANDRILL_API_KEY`: define mandrill provider api key.

`PROVIDER_SES_AWS_KEY_ID`: define aws key id.

`PROVIDER_SES_AWS_SECRET_KEY`: define aws secret key.

`PROVIDER_SES_AWS_REGION`: define aws region.

## Commands (Development)

`make build`: build restaurants service for osx.

`make linux`: build restaurants service for linux os.

`make docker .`: build docker.

`make compose`: start docker-docker.

`make stop`: stop docker-docker.

`make run`: run email service.

`docker run -it -p 5010:5010 email-api`: run docker.
