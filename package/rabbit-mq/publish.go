package rabbitmq

import "github.com/streadway/amqp"


func (r *RabbitMQClient) Publish(queueName string, body []byte) error {
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
