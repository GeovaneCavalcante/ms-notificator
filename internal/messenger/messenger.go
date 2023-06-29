package messenger

type MessengerConfiguration struct {
	TopicName     string
	Configuration map[string]interface{}
}

type MessageResponse struct {
	ID string
}

type Messenger interface {
	PublishMessage(message string) (*MessageResponse, error)
}
