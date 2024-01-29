package actorsystem

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

// Actor represents an actor that connects to AMQP, creates/binds a queue, and begins consuming messages.
type Actor struct {
	Exchange   string
	RoutingKey string
	Queue      string
	Handler    func(message []byte) error
}

// Start initializes the actor, connects to AMQP, creates/binds a queue, and starts consuming messages.
func (a *Actor) Start(amqpURI string) error {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("failed to connect to AMQP: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Use "default" exchange if not specified
	exchange := a.Exchange
	if exchange == "" {
		exchange = ""
	}

	// Use queue name as routing key if not specified
	routingKey := a.RoutingKey
	if routingKey == "" {
		routingKey = a.Queue
	}

	err = ch.ExchangeDeclare(
		exchange, // exchange name
		"direct", // exchange type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = ch.QueueDeclare(
		a.Queue, // queue name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		a.Queue,    // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		a.Queue, // queue name
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case msg := <-msgs:
				err := a.Handler(msg.Body)
				if err != nil {
					log.Printf("Error processing message: %v", err)
				} else {
					msg.Ack(false)
				}
			case <-sig:
				log.Println("Received interrupt. Shutting down...")
				os.Exit(0)
			}
		}
	}()

	log.Println("Actor system started. Waiting for messages...")
	select {}
}
