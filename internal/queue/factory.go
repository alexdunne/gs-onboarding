package queue

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Message represents the structure of the messages being sent
type Message struct {
	ID int `json:"id"`
}

// Queue is a interface to expose methods to interact with a queue
type Queue interface {
	Publish(msg *Message) error
	Consume(ctx context.Context) (<-chan *Message, error)
}

type client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	logger  *zap.Logger
}

// New creates a connection to a RMQ instance and configures the necessary queues
func New(connStr string, queueName string, logger *zap.Logger) (*client, error) {
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to rabbitmq")
	}

	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a channel")
	}

	q, err := amqpChan.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare queue")
	}

	c := &client{
		conn:    conn,
		channel: amqpChan,
		queue:   q,
		logger:  logger,
	}

	return c, nil
}

// Close closes any created channels and connections
func (c *client) Close() error {
	if err := c.channel.Close(); err != nil {
		return errors.Wrap(err, "failed to close the channel")
	}

	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "failed to close the connection")
	}

	return nil
}

// Publish sends a message to a queue
func (c *client) Publish(msg *Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	return c.channel.Publish(
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

// Consumer continuously receives messages from a queue and sends them to a returned channel
func (c *client) Consume(ctx context.Context) (<-chan *Message, error) {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	messages := make(chan *Message)

	go func() {
		defer close(messages)

		c.logger.Info("consuming messages")

		for {
			select {
			case <-ctx.Done():
				return
			case in, ok := <-msgs:
				if !ok {
					return
				}

				var msg *Message
				if err := json.Unmarshal(in.Body, &msg); err != nil {
					c.logger.Info("failed to convert incoming message Message struct")
					in.Nack(false, true)
					continue
				}

				messages <- msg
			}
		}
	}()

	return messages, nil
}
