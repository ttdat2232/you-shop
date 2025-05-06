package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/TechwizsonORG/image-service/config/model"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Service struct {
	config model.RabbitMqConfig
	logger zerolog.Logger
}

func NewRpcService(config model.RabbitMqConfig, logger zerolog.Logger) *Service {
	logger = logger.With().
		Str("Infrastructure", "RPC Service").
		Logger()
	return &Service{
		config: config,
		logger: logger,
	}
}

func fatalOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (s *Service) NewRpcQueue(rpcQueueName string, handle func(data string) string) {
	conn, err := amqp091.Dial(s.config.GetAmqpServerUrl())
	fatalOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	fatalOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(fmt.Sprintf("%s.%s", rpcQueueName, "rpc.response"), false, false, false, false, nil)

	fatalOnError(err, "Failed to declare a queue")

	err = ch.Qos(1, 0, false)
	fatalOnError(err, "Failed to set QoS")

	_, err = ch.QueueDeclare(fmt.Sprintf("%s.%s", rpcQueueName, "rpc.request"), false, false, false, false, nil)
	if err != nil {
		fatalOnError(err, "Failed to declare a listening request queue")
	}

	msgs, err := ch.Consume(
		fmt.Sprintf("%s.%s", rpcQueueName, "rpc.request"),
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	fatalOnError(err, "Failed to register a consumer")

	var forever chan struct{}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for d := range msgs {
			res := handle(string(d.Body))
			var js json.RawMessage
			if err = json.Unmarshal([]byte(res), &js); err != nil {
				fatalOnError(err, "Failed to unmarshal json")
			}
			err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(res),
			})
			fatalOnError(err, "Failed to publish a message")
			d.Ack(false)
		}
	}()

	<-forever
}

func (s *Service) Req(rpcQueueName string, request string) string {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    conn, err := amqp091.Dial(s.config.GetAmqpServerUrl())
    fatalOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    fatalOnError(err, "Failed to open a channel")
    defer ch.Close()

    // Declare request queue
    requestQueue := fmt.Sprintf("%s.%s", rpcQueueName, "rpc.request")
    _, err = ch.QueueDeclare(
        requestQueue,
        false, // durable
        false, // auto-delete
        false, // exclusive
        false, // no-wait
        nil,
    )
    fatalOnError(err, "Failed to declare request queue")
	responseQueueName := fmt.Sprintf("%s.%s", rpcQueueName, "rpc.response")
    msgs, err := ch.Consume(
        responseQueueName,
        "",    // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,
    )
    fatalOnError(err, "Failed to register a consumer")

    corrId := randomString(32)

    // Validate JSON request
    if !json.Valid([]byte(request)) {
        fatalOnError(fmt.Errorf("invalid JSON"), "Invalid request format")
        return ""
    }

    // Publish request with context
    err = ch.PublishWithContext(ctx,
        "",          // exchange
        requestQueue, // routing key
        false,       // mandatory
        false,       // immediate
        amqp091.Publishing{
            ContentType:   "application/json",
            CorrelationId: corrId,
            ReplyTo:       responseQueueName,
            Body:          []byte(request),
        })
    if err != nil {
        s.logger.Error().Err(err).Msg("Failed to publish request")
        return ""
    }

    // Handle response or timeout
    for {
        select {
        case d, ok := <-msgs:
            if !ok {
                s.logger.Error().Msg("Response channel closed")
                return ""
            }

            if d.CorrelationId == corrId {
                if err := d.Ack(false); err != nil {
                    s.logger.Error().Err(err).Msg("Failed to ack response")
                }
                return string(d.Body)
            } else {
                // Requeue unexpected messages
                if err := d.Nack(false, true); err != nil {
                    s.logger.Error().Err(err).Msg("Failed to requeue message")
                }
            }
        case <-ctx.Done():
            s.logger.Error().Err(ctx.Err()).Msg("Request timed out")
            return ""
        }
    }

}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
