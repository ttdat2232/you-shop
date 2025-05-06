package rabbitmq

import (
	"encoding/json"

	"github.com/TechwizsonORG/order-service/config/model"
	"github.com/TechwizsonORG/order-service/usecase/messagequeue"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type DefaultMessageQueue struct {
	rabbitMqConfig model.RabbitMqConfig
	logger         zerolog.Logger
}

func NewDefaultMessageQueue(rabbitMq model.RabbitMqConfig, logger zerolog.Logger) *DefaultMessageQueue {
	logger = logger.
		With().
		Str("infrastructure", "rabbitMQ").
		Logger()

	return &DefaultMessageQueue{
		rabbitMqConfig: rabbitMq,
		logger:         logger,
	}
}

func (mq *DefaultMessageQueue) failOnError(err error, msg string) {
	if err != nil {
		mq.logger.Error().Err(err).Msg(msg)
	}
}
func (mq *DefaultMessageQueue) Publish(exchangeConfig messagequeue.ExchangeConfig, queueConfig messagequeue.QueueConfig, data any) {
	conn, err := amqp091.Dial(mq.rabbitMqConfig.GetAmqpServerUrl())
	mq.failOnError(err, "failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	mq.failOnError(err, "failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(exchangeConfig.ExchangeName, string(exchangeConfig.Type), exchangeConfig.Durable, exchangeConfig.AutoDelete, exchangeConfig.Internal, exchangeConfig.NoWait, nil)
	mq.failOnError(err, "failed to declare exchange")

	jsonData, err := json.Marshal(data)
	mq.failOnError(err, "failed to convert data into json")

	err = ch.Publish(exchangeConfig.ExchangeName, queueConfig.RoutingKey, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        jsonData,
	})
	mq.failOnError(err, "failed to publish message")
	mq.logger.Debug().Msg("Published message successfully")
	mq.logger.Trace().Msgf("value %s", jsonData)

}
func (mq *DefaultMessageQueue) Consume(exchangeConfig messagequeue.ExchangeConfig, queueConfig messagequeue.QueueConfig, handler func(data string) error) {
	conn, err := amqp091.Dial(mq.rabbitMqConfig.GetAmqpServerUrl())
	mq.failOnError(err, "failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	mq.failOnError(err, "failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(exchangeConfig.ExchangeName, string(exchangeConfig.Type), exchangeConfig.Durable, exchangeConfig.AutoDelete, exchangeConfig.Internal, exchangeConfig.NoWait, nil)
	mq.failOnError(err, "failed to declare exchange")

	q, err := ch.QueueDeclare(queueConfig.QueueName, queueConfig.Durable, queueConfig.DeleteUnused, queueConfig.Exclusive, queueConfig.NoWait, nil)
	mq.failOnError(err, "failed to declare queue")

	err = ch.QueueBind(q.Name, queueConfig.RoutingKey, exchangeConfig.ExchangeName, queueConfig.NoWait, nil)
	mq.failOnError(err, "Failed to bind queue")

	msgs, err := ch.Consume(queueConfig.QueueName, "", true, queueConfig.Exclusive, false, queueConfig.NoWait, nil)
	mq.failOnError(err, "failed to register a consumer")

	var forever chan interface{}
	go func() {
		for d := range msgs {
			err = handler(string(d.Body))
			if err != nil {
				mq.logger.Error().Err(err).Msg("")
			} else {
				d.Ack(false)
			}
		}
	}()
	<-forever

}
