package server

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gateway-service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

const accessLogPublishTimeout = 2 * time.Second
const accessLogQueueSize = 1024

type AccessLogPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	log     *log.Helper
	bodyCh  chan []byte
	doneCh  chan struct{}
}

func NewAccessLogPublisher(c *conf.RabbitMQConfig, logger log.Logger) (*AccessLogPublisher, func(), error) {
	helper := log.NewHelper(logger)
	if c == nil || c.Addr == "" {
		return nil, func() {}, errors.New("rabbitmq addr is required")
	}

	queue := c.Topic
	if queue == "" {
		queue = "access_log"
	}

	conn, err := amqp.Dial(c.Addr)
	if err != nil {
		return nil, nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	if _, err = channel.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, nil, err
	}

	publisher := &AccessLogPublisher{
		conn:    conn,
		channel: channel,
		queue:   queue,
		log:     helper,
		bodyCh:  make(chan []byte, accessLogQueueSize),
		doneCh:  make(chan struct{}),
	}
	go publisher.run()

	cleanup := func() {
		helper.Info("closing access log rabbitmq publisher")
		close(publisher.bodyCh)
		<-publisher.doneCh
		if err := channel.Close(); err != nil {
			helper.Errorf("close rabbitmq channel failed: %v", err)
		}
		if err := conn.Close(); err != nil {
			helper.Errorf("close rabbitmq connection failed: %v", err)
		}
	}
	return publisher, cleanup, nil
}

func (p *AccessLogPublisher) Publish(data networkRequest) {
	if p == nil {
		return
	}
	body, err := json.Marshal(data)
	if err != nil {
		p.log.Errorf("marshal access log failed: %v", err)
		return
	}

	select {
	case p.bodyCh <- body:
	default:
		p.log.Errorf("access log rabbitmq queue is full, drop request_id=%s", data.RequestId)
	}
}

func (p *AccessLogPublisher) run() {
	defer close(p.doneCh)

	for body := range p.bodyCh {
		p.publish(body)
	}
}

func (p *AccessLogPublisher) publish(body []byte) {
	publishCtx, cancel := context.WithTimeout(context.Background(), accessLogPublishTimeout)
	defer cancel()

	if err := p.channel.PublishWithContext(
		publishCtx,
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	); err != nil {
		p.log.Errorf("publish access log to rabbitmq failed: %v", err)
	}
}
