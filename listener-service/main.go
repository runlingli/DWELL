package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Connected to RabbitMQ! Starting listeners...")

	// Create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// Start listener for log events in a goroutine
	go func() {
		log.Println("Starting log events listener...")
		err := consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
		if err != nil {
			log.Printf("Error in log listener: %v", err)
		}
	}()

	// Start listener for app events (mail, notifications) in a goroutine
	go func() {
		log.Println("Starting app events listener...")
		// Need a new consumer for the second listener
		appConsumer, err := event.NewConsumer(rabbitConn)
		if err != nil {
			log.Printf("Error creating app consumer: %v", err)
			return
		}
		err = appConsumer.ListenForAppEvents()
		if err != nil {
			log.Printf("Error in app events listener: %v", err)
		}
	}()

	log.Println("Listener service is running. Press CTRL+C to exit.")

	// Block forever
	forever := make(chan bool)
	<-forever
}

// connect attempts to connect to RabbitMQ with exponential backoff
func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
