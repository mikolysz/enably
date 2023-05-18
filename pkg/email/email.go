package email

type Message struct {
	Recipient        string
	Subject          string
	PlainTextContent string
	HTMLContent      string
}

type Sender interface {
	Send(m Message) error
}
