package email

import (
	"github.com/microapis/users-api"
)

const (
	// Channel ...
	Channel = "email"
)

// Message ...
type Message struct {
	From     string   `json:"from"`
	FromName string   `json:"from_name"`
	To       []string `json:"to"`
	ReplyTo  []string `json:"reply_to"`
	Subject  string   `json:"subject"`
	Text     string   `json:"text"`
	HTML     string   `json:"html"`
	Provider string   `json:"provider"`
	Status   string   `json:"status"`
}

// MailingTemplates ...
type MailingTemplates struct {
	Signup          func(u *users.User) error
	VerifyEmail     func(u *users.User, token string) error
	ForgotPassword  func(u *users.User, token string) error
	PasswordChanged func(u *users.User) error
}
