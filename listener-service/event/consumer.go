package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer handles consuming messages from RabbitMQ
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer creates a new consumer instance
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup initializes the consumer by declaring all exchanges
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return DeclareAllExchanges(channel)
}

// Payload represents a generic message payload
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// MailPayload represents an email message
type MailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// VerificationMailPayload represents a verification email
type VerificationMailPayload struct {
	To               string `json:"to"`
	FirstName        string `json:"first_name"`
	VerificationCode string `json:"verification_code"`
	Type             string `json:"type"`
}

// Listen starts listening for messages on specified topics
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Create queue for logs exchange
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// Bind topics to logs exchange
	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			ExchangeLogs, // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [%s, %s]\n", ExchangeLogs, q.Name)
	<-forever

	return nil
}

// ListenForAppEvents listens for application events (mail, notifications, etc.)
func (consumer *Consumer) ListenForAppEvents() error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Create a durable queue for app events
	q, err := declareDurableQueue(ch, "app_events_queue")
	if err != nil {
		return err
	}

	// Bind mail topics
	mailTopics := []string{"mail.*", "notification.*"}
	for _, topic := range mailTopics {
		err = ch.QueueBind(
			q.Name,      // queue name
			topic,       // routing key pattern
			ExchangeApp, // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	// Consume with manual acknowledgment for reliability
	messages, err := ch.Consume(
		q.Name,
		"",
		false, // auto-ack = false for manual acknowledgment
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received app event with routing key: %s", d.RoutingKey)

			var err error
			switch d.RoutingKey {
			case RoutingMailSend:
				err = handleMailEvent(d.Body)
			case RoutingMailVerification:
				err = handleVerificationMailEvent(d.Body)
			case RoutingMailPasswordReset:
				err = handleVerificationMailEvent(d.Body)
			default:
				log.Printf("Unknown routing key: %s", d.RoutingKey)
			}

			if err != nil {
				log.Printf("Error handling event: %v", err)
				// Negative acknowledgment - requeue the message
				d.Nack(false, true)
			} else {
				// Positive acknowledgment
				d.Ack(false)
			}
		}
	}()

	fmt.Printf("Waiting for app events [Exchange, Queue] [%s, %s]\n", ExchangeApp, q.Name)
	<-forever

	return nil
}

// handlePayload processes log events
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// Handle auth events if needed
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// handleMailEvent processes mail send events
func handleMailEvent(body []byte) error {
	var mail MailPayload
	if err := json.Unmarshal(body, &mail); err != nil {
		return err
	}

	log.Printf("Sending mail to: %s, subject: %s", mail.To, mail.Subject)
	return sendMailToService(mail)
}

// handleVerificationMailEvent processes verification email events
func handleVerificationMailEvent(body []byte) error {
	var mail VerificationMailPayload
	if err := json.Unmarshal(body, &mail); err != nil {
		return err
	}

	log.Printf("Sending verification mail to: %s, type: %s", mail.To, mail.Type)

	// Convert to regular mail format
	subject := "Email Verification"
	if mail.Type == "password_reset" {
		subject = "Password Reset Code"
	}

	message := fmt.Sprintf("Hello %s,\n\nYour verification code is: %s\n\nThis code will expire in 10 minutes.",
		mail.FirstName, mail.VerificationCode)

	return sendMailToService(MailPayload{
		To:      mail.To,
		Subject: subject,
		Message: message,
	})
}

// sendMailToService sends the mail via HTTP to mailer-service
func sendMailToService(mail MailPayload) error {
	jsonData, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	mailServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("mailer service returned status: %d", response.StatusCode)
	}

	log.Printf("Mail sent successfully to: %s", mail.To)
	return nil
}

// logEvent sends log to logger service
func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("logger service returned status: %d", response.StatusCode)
	}

	return nil
}
