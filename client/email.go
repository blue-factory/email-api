package client

import (
	"context"
	"encoding/json"
	"errors"

	email "github.com/microapis/email-api"
	"github.com/microapis/messages-core/message"
	"github.com/microapis/messages-core/proto"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	Client proto.SchedulerServiceClient
}

// New ...
func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := proto.NewSchedulerServiceClient(conn)

	return &Client{
		Client: c,
	}, nil
}

// Send ...
func (c *Client) Send(e *email.Message, delay int64) (string, error) {
	// validates email existence
	if e == nil {
		return "", errors.New("invalid email")
	}

	// validates email params
	// TODO(ca): check how verify e.Text, e.HTML, e.ReplyTo and e.Status
	if e.From == "" || e.FromName == "" || e.To[0] == "" || e.Subject == "" || e.Provider == "" {
		return "", errors.New("invalid email params")
	}

	// validates delay param
	if delay < 0 {
		return "", errors.New("invalid delay")
	}

	// pase email struct
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}

	e.Status = message.Pending

	gr, err := c.Client.Put(context.Background(), &proto.MessagePutRequest{
		Channel:  email.Channel,
		Provider: e.Provider,
		Content:  string(b),
		Delay:    delay,
	})
	if err != nil {
		return "", err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return "", errors.New(msg)
	}

	data := gr.GetData()
	id := data.GetId()

	return id, nil
}

// Get ...
func (c *Client) Get(ID string) (*email.Message, error) {
	// validates id param
	if ID == "" {
		return nil, errors.New("invalid ID")
	}

	gr, err := c.Client.Get(context.Background(), &proto.MessageGetRequest{
		Id: ID,
	})
	if err != nil {
		return nil, err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return nil, errors.New(msg)
	}

	data := gr.GetData()
	provider := data.GetProvider()
	content := data.GetContent()
	status := data.GetStatus()

	m := &email.Message{}
	err = json.Unmarshal([]byte(content), m)
	if err != nil {
		return nil, err
	}

	m.Provider = provider
	m.Status = status

	return m, nil
}

// Update ...
func (c *Client) Update(ID string, e *email.Message) error {
	// validates email existence
	if e == nil {
		return errors.New("invalid email")
	}

	// validates email params
	// TODO(ca): check how verify e.Text, e.HTML, e.ReplyTo and e.Status
	if e.From == "" || e.FromName == "" || e.To[0] == "" || e.Subject == "" || e.Provider == "" {
		return errors.New("invalid email params")
	}

	// pase email struct
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	gr, err := c.Client.Update(context.Background(), &proto.MessageUpdateRequest{
		Id:      ID,
		Content: string(b),
	})
	if err != nil {
		return err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return errors.New(msg)
	}

	return nil
}

// Cancel ...
func (c *Client) Cancel(ID string) error {
	// validates id param
	if ID == "" {
		return errors.New("invalid ID")
	}

	gr, err := c.Client.Cancel(context.Background(), &proto.MessageCancelRequest{
		Id: ID,
	})
	if err != nil {
		return err
	}

	msg := gr.GetError().GetMessage()
	if msg != "" {
		return errors.New(msg)
	}

	if err != nil {
		return err
	}

	return nil
}
