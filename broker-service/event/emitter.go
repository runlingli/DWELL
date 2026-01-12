package event

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter handles publishing events to RabbitMQ
type Emitter struct {
	connection *amqp.Connection
}

// MailEvent represents an email to be sent
type MailEvent struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// VerificationMailEvent represents a verification email
type VerificationMailEvent struct {
	To               string `json:"to"`
	FirstName        string `json:"first_name"`
	VerificationCode string `json:"verification_code"`
	Type             string `json:"type"` // "signup" or "password_reset"
}

// NotificationEvent represents a notification
type NotificationEvent struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

// setup initializes the emitter by declaring all exchanges
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return DeclareAllExchanges(channel)
}

// Push sends a message to the logs exchange (legacy method)
func (e *Emitter) Push(event string, severity string) error {
	return e.pushToExchange(ExchangeLogs, severity, []byte(event))
}

// PushToApp sends a message to the app events exchange
func (e *Emitter) PushToApp(routingKey string, payload any) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return e.pushToExchange(ExchangeApp, routingKey, jsonData)
}

// pushToExchange is the internal method that publishes to any exchange
func (e *Emitter) pushToExchange(exchange, routingKey string, body []byte) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Printf("Pushing to exchange [%s] with routing key [%s]", exchange, routingKey)

	err = channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Message survives broker restart
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// SendMail publishes a mail event to RabbitMQ
func (e *Emitter) SendMail(to, subject, message string) error {
	event := MailEvent{
		To:      to,
		Subject: subject,
		Message: message,
	}
	return e.PushToApp(RoutingMailSend, event)
}

// SendVerificationMail publishes a verification email event
func (e *Emitter) SendVerificationMail(to, firstName, code, mailType string) error {
	event := VerificationMailEvent{
		To:               to,
		FirstName:        firstName,
		VerificationCode: code,
		Type:             mailType,
	}
	return e.PushToApp(RoutingMailVerification, event)
}

// SendNotification publishes a notification event
func (e *Emitter) SendNotification(userID int, title, message, notifType string) error {
	event := NotificationEvent{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
	}
	return e.PushToApp(RoutingNotification, event)
}

// NewEventEmitter creates a new Emitter instance
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
