package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"platform/internal/config"
	"platform/internal/models"
	"platform/pkg/logger"
)

type Publisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewPublisher(cfg *config.Config) (*Publisher, error) {
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

	return &Publisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Publisher) PublishRequest(request *models.Request) error {
	// Convert request to JSON
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Publish message
	err = p.channel.PublishWithContext(
		nil,           // context
		"",            // exchange
		p.queue.Name,  // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	logger.Info("Published request to queue", "request_id", request.ID)
	return nil
} 