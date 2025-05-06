package model

import "fmt"

type RabbitMqConfig struct {
	Host     string
	Username string
	Password string
	Port     int
	Vhost	string
}

func (r *RabbitMqConfig) GetAmqpServerUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.Username, r.Password, r.Host, r.Port, r.Vhost)
}
