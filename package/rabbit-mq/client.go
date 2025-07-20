package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/dika22/news-service/package/config"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *config.Config
}

type IRabbitMQClient interface {
	Publish(queue string, body []byte) error
	Consume(queueName string, handler func(amqp.Delivery)) error 
	Close()
}

func NewRabbitMQClient(cfg *config.Config) (*RabbitMQClient, error) {
	client := &RabbitMQClient{cfg: cfg}
	if err := client.connect(); err != nil {
		return nil, err
	}
	return client, nil
}

func (r *RabbitMQClient) connect() error {
	var err error
    url := fmt.Sprintf("amqp://%v:%v@%v:%v/", r.cfg.MessageBrokerUsername, r.cfg.MessageBrokerPassword, r.cfg.MessageBrokerURL, r.cfg.MessageBrokerPort)
	for i := 0; i < 5; i++ {
		r.conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Retrying RabbitMQ connection... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return err
	}

	r.channel, err = r.conn.Channel()
	return err
}

func (r *RabbitMQClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
