package messagesemail

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
