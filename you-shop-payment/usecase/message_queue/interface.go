package messagequeue

type ExchangeType string

const (
	Fanout  ExchangeType = "fanout"
	Direct  ExchangeType = "direct"
	Topic   ExchangeType = "topic"
	Headers ExchangeType = "headers"
)

type ExchangeConfig struct {
	ExchangeName string
	Type         ExchangeType
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
}

type QueueConfig struct {
	QueueName    string
	RoutingKey   string
	Durable      bool
	DeleteUnused bool
	Exclusive    bool
	NoWait       bool
}

func NewDefaultQueueConfig(name string, routingKey string) *QueueConfig {
	return &QueueConfig{
		QueueName:    name,
		RoutingKey:   routingKey,
		Durable:      false,
		DeleteUnused: false,
		Exclusive:    false,
		NoWait:       false,
	}
}

func NewDefaultExchangeConfig(name string, exchangeType ExchangeType) *ExchangeConfig {
	return &ExchangeConfig{
		ExchangeName: name,
		Type:         exchangeType,
		Durable:      false,
		AutoDelete:   false,
		Internal:     false,
		NoWait:       false,
	}
}

type MessageQueue interface {

	// Data will be tried to parse in JSON form
	Publish(exchangeConfig ExchangeConfig, queueConfig QueueConfig, data any)

	// Data parameter in handler will be a string in JSON form
	Consume(exchangeConfig ExchangeConfig, queueConfig QueueConfig, handler func(data string) error)
}
