package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// Exchange names
const (
	ExchangeLogs = "logs_topic" // For logging events
	ExchangeApp  = "app_events" // For application events (mail, notifications, etc.)
)

// Routing keys for app events
const (
	RoutingMailSend         = "mail.send"
	RoutingMailVerification = "mail.verification"
	RoutingMailPasswordReset = "mail.password_reset"
	RoutingNotification     = "notification.send"
)

// declareExchange declares the logs exchange (legacy)
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		ExchangeLogs, // name
		"topic",      // type
		true,         // durable?
		false,        // auto-deleted?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguments?
	)
}

// declareAppExchange declares the app events exchange
func declareAppExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		ExchangeApp, // name
		"topic",     // type
		true,        // durable?
		false,       // auto-deleted?
		false,       // internal?
		false,       // no-wait?
		nil,         // arguments?
	)
}

// DeclareAllExchanges declares all exchanges needed by the application
func DeclareAllExchanges(ch *amqp.Channel) error {
	if err := declareExchange(ch); err != nil {
		return err
	}
	return declareAppExchange(ch)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unused?
		true,  // exclusive?
		false, // no-wait?
		nil,   // arguments?
	)
}

// declareDurableQueue declares a durable queue for reliable message delivery
func declareDurableQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,  // name
		true,  // durable - survives broker restart
		false, // delete when unused?
		false, // exclusive?
		false, // no-wait?
		nil,   // arguments?
	)
}
