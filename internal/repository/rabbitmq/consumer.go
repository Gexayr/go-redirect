package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"platform/internal/config"
	"platform/internal/models"
	"platform/pkg/logger"
)

type Consumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	// Create RabbitMQ connection URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	// Connect to RabbitMQ
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.Queue, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) Consume(handler func(*models.Request) error) error {
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var request models.Request
			if err := json.Unmarshal(d.Body, &request); err != nil {
				logger.Error("Failed to unmarshal request", "error", err.Error())
				// Reject the message and requeue it
				d.Reject(true)
				continue
			}

			if err := handler(&request); err != nil {
				logger.Error("Failed to process request", "error", err.Error())
				// Reject the message and requeue it
				d.Reject(true)
				continue
			}

			// Acknowledge the message after successful processing
			d.Ack(false)
			logger.Info("Successfully processed request", "request_id", request.ID)
		}
	}()

	return nil
} 