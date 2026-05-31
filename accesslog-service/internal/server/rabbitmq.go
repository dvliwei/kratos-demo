package server

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConfig struct {
	URL           string `json:"url"`
	Queue         string `json:"queue"`
	ConsumerTag   string `json:"consumer_tag"`
	AutoAck       bool   `json:"auto_ack"`
	PrefetchCount uint32 `json:"prefetch_count"`
}

func (c *RabbitMQConfig) GetURL() string {
	if c != nil {
		return c.URL
	}
	return ""
}

func (c *RabbitMQConfig) GetQueue() string {
	if c != nil {
		return c.Queue
	}
	return ""
}

func (c *RabbitMQConfig) GetConsumerTag() string {
	if c != nil {
		return c.ConsumerTag
	}
	return ""
}

func (c *RabbitMQConfig) GetAutoAck() bool {
	if c != nil {
		return c.AutoAck
	}
	return false
}

func (c *RabbitMQConfig) GetPrefetchCount() uint32 {
	if c != nil {
		return c.PrefetchCount
	}
	return 0
}

// RabbitMQConsumer consumes access log messages and prints them for now.
type RabbitMQConsumer struct {
	conf   *RabbitMQConfig
	log    *log.Helper
	conn   *amqp.Connection
	ch     *amqp.Channel
	cancel context.CancelFunc
	done   chan struct{}
}

func NewRabbitMQConsumer(c *RabbitMQConfig, logger log.Logger) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		conf: c,
		log:  log.NewHelper(log.With(logger, "module", "server/rabbitmq")),
		done: make(chan struct{}),
	}
}

func (r *RabbitMQConsumer) Start(ctx context.Context) error {
	if r.conf == nil || r.conf.GetURL() == "" || r.conf.GetQueue() == "" {
		r.log.Warn("rabbitmq consumer disabled: missing url or queue")
		close(r.done)
		return nil
	}

	conn, err := amqp.Dial(r.conf.GetURL())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	if prefetch := int(r.conf.GetPrefetchCount()); prefetch > 0 {
		if err := ch.Qos(prefetch, 0, false); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return err
		}
	}

	if _, err := ch.QueueDeclare(r.conf.GetQueue(), true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	deliveries, err := ch.Consume(
		r.conf.GetQueue(),
		r.conf.GetConsumerTag(),
		r.conf.GetAutoAck(),
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	runCtx, cancel := context.WithCancel(ctx)
	r.conn = conn
	r.ch = ch
	r.cancel = cancel
	r.done = make(chan struct{})

	go r.consume(runCtx, deliveries)
	r.log.Infof("rabbitmq consumer started: queue=%s auto_ack=%t", r.conf.GetQueue(), r.conf.GetAutoAck())
	return nil
}

func (r *RabbitMQConsumer) Stop(context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}

	if r.ch != nil {
		_ = r.ch.Cancel(r.conf.GetConsumerTag(), false)
		_ = r.ch.Close()
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}

	<-r.done
	r.log.Info("rabbitmq consumer stopped")
	return nil
}

func (r *RabbitMQConsumer) consume(ctx context.Context, deliveries <-chan amqp.Delivery) {
	defer close(r.done)

	for {
		select {
		case <-ctx.Done():
			return
		case delivery, ok := <-deliveries:
			if !ok {
				return
			}
			r.log.Infof("received rabbitmq message: exchange=%s routing_key=%s body=%s", delivery.Exchange, delivery.RoutingKey, string(delivery.Body))
			if !r.conf.GetAutoAck() {
				if err := delivery.Ack(false); err != nil && !errors.Is(err, amqp.ErrClosed) {
					r.log.Errorf("ack rabbitmq message failed: %v", err)
				}
			}
		}
	}
}
