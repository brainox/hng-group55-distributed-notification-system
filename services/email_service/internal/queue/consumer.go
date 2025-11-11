package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/idempotency"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/models"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/retry"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/sender"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/internal/template"
	"github.com/brainox/hng-group55-distributed-notification-system/services/email_service/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

type Consumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	queueName      string
	workerCount    int
	templateClient *template.Client
	emailSender    sender.EmailSender
	publisher      *Publisher
	idempotency    *idempotency.Checker
	retryHandler   *retry.Handler
	circuitBreaker *gobreaker.CircuitBreaker
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

type ConsumerConfig struct {
	URL            string
	QueueName      string
	WorkerCount    int
	TemplateClient *template.Client
	EmailSender    sender.EmailSender
	Publisher      *Publisher
	Idempotency    *idempotency.Checker
	RetryHandler   *retry.Handler
	CircuitBreaker *gobreaker.CircuitBreaker
}

func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS to process one message at a time per worker
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	// Declare queue
	_, err = channel.QueueDeclare(
		cfg.QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		conn:           conn,
		channel:        channel,
		queueName:      cfg.QueueName,
		workerCount:    cfg.WorkerCount,
		templateClient: cfg.TemplateClient,
		emailSender:    cfg.EmailSender,
		publisher:      cfg.Publisher,
		idempotency:    cfg.Idempotency,
		retryHandler:   cfg.RetryHandler,
		circuitBreaker: cfg.CircuitBreaker,
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",    // consumer tag
		false, // auto-ack (we want manual ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	logger.Log.Info("starting email consumer",
		zap.String("queue", c.queueName),
		zap.Int("workers", c.workerCount),
	)

	// Start worker pool
	for i := 0; i < c.workerCount; i++ {
		c.wg.Add(1)
		go c.worker(i, msgs)
	}

	return nil
}

func (c *Consumer) worker(id int, msgs <-chan amqp.Delivery) {
	defer c.wg.Done()

	logger.Log.Info("worker started", zap.Int("worker_id", id))

	for {
		select {
		case <-c.ctx.Done():
			logger.Log.Info("worker stopping", zap.Int("worker_id", id))
			return
		case msg, ok := <-msgs:
			if !ok {
				logger.Log.Info("message channel closed", zap.Int("worker_id", id))
				return
			}
			c.processMessage(msg)
		}
	}
}

func (c *Consumer) processMessage(delivery amqp.Delivery) {
	var emailMsg models.EmailMessage
	if err := json.Unmarshal(delivery.Body, &emailMsg); err != nil {
		logger.Log.Error("failed to unmarshal message", zap.Error(err))
		delivery.Nack(false, false) // Don't requeue invalid messages
		return
	}

	logger.Log.Info("processing email",
		zap.String("id", emailMsg.ID),
		zap.String("correlation_id", emailMsg.CorrelationID),
		zap.String("recipient", emailMsg.Recipient),
		zap.String("template_id", emailMsg.TemplateID),
	)

	// Check idempotency
	processed, err := c.idempotency.IsProcessed(c.ctx, emailMsg.ID)
	if err != nil {
		logger.Log.Error("failed to check idempotency", zap.Error(err))
	}
	if processed {
		logger.Log.Info("message already processed", zap.String("id", emailMsg.ID))
		delivery.Ack(false)
		return
	}

	// Process with retry
	err = c.processWithRetry(&emailMsg)

	if err != nil {
		logger.Log.Error("failed to process email after retries",
			zap.Error(err),
			zap.String("id", emailMsg.ID),
		)

		// Publish failed status
		c.publishStatus(emailMsg.ID, emailMsg.CorrelationID, "failed", err.Error())

		// Don't requeue - message goes to DLQ or is discarded
		delivery.Nack(false, false)
		return
	}

	// Mark as processed
	if err := c.idempotency.MarkProcessed(c.ctx, emailMsg.ID); err != nil {
		logger.Log.Error("failed to mark as processed", zap.Error(err))
	}

	// Publish success status
	c.publishStatus(emailMsg.ID, emailMsg.CorrelationID, "sent", "")

	// Acknowledge message
	delivery.Ack(false)

	logger.Log.Info("email sent successfully",
		zap.String("id", emailMsg.ID),
		zap.String("recipient", emailMsg.Recipient),
	)
}

func (c *Consumer) processWithRetry(emailMsg *models.EmailMessage) error {
	var lastErr error
	maxAttempts := 5

	for attempt := 0; attempt <= maxAttempts; attempt++ {
		if attempt > 0 {
			c.retryHandler.Wait(attempt-1, emailMsg.CorrelationID)
		}

		err := c.processEmail(emailMsg)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		if !c.retryHandler.ShouldRetry(err, attempt) {
			logger.Log.Warn("not retrying",
				zap.Error(err),
				zap.Int("attempt", attempt),
				zap.String("correlation_id", emailMsg.CorrelationID),
			)
			break
		}

		logger.Log.Warn("retrying after error",
			zap.Error(err),
			zap.Int("attempt", attempt),
			zap.String("correlation_id", emailMsg.CorrelationID),
		)
	}

	return lastErr
}

func (c *Consumer) processEmail(emailMsg *models.EmailMessage) error {
	// Fetch template
	tmpl, err := c.templateClient.FetchTemplate(c.ctx, emailMsg.TemplateID)
	if err != nil {
		return fmt.Errorf("failed to fetch template: %w", err)
	}

	// Render subject
	subject, err := template.RenderTemplate(tmpl.Subject, emailMsg.Variables)
	if err != nil {
		return fmt.Errorf("failed to render subject: %w", err)
	}

	// Render body
	body, err := template.RenderTemplate(tmpl.Body, emailMsg.Variables)
	if err != nil {
		return fmt.Errorf("failed to render body: %w", err)
	}

	// Send email with circuit breaker
	_, err = c.circuitBreaker.Execute(func() (interface{}, error) {
		return nil, c.emailSender.Send(emailMsg.Recipient, subject, body)
	})

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (c *Consumer) publishStatus(notificationID, correlationID, status, errorMsg string) {
	statusMsg := models.StatusMessage{
		NotificationID: notificationID,
		CorrelationID:  correlationID,
		Status:         status,
		Timestamp:      time.Now(),
		Error:          errorMsg,
		Provider:       c.emailSender.GetProviderName(),
	}

	if err := c.publisher.PublishStatus(c.ctx, statusMsg); err != nil {
		logger.Log.Error("failed to publish status", zap.Error(err))
	}
}

func (c *Consumer) Stop() {
	logger.Log.Info("stopping consumer...")
	c.cancel()
	c.wg.Wait()

	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}

	logger.Log.Info("consumer stopped")
}
