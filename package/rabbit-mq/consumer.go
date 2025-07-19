package rabbitmq

import (
	"github.com/streadway/amqp"
)

func (r *RabbitMQClient) Consume(queueName string, handler func(amqp.Delivery)) error {
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		handler(msg)
	}

	return nil
}
