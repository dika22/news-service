package task

import (
	"context"
	"encoding/json"
	"fmt"
	"news-service/package/config"
	"news-service/package/connection/elasticsearch"
	rabbitmq "news-service/package/rabbit-mq"
	"news-service/package/structs"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

type TaskWorker struct{
	conf     *config.Config
	mqClient *rabbitmq.RabbitMQClient
	esClient *elasticsearch.Elasticsearch
}

func (t *TaskWorker) Run(ctx context.Context) error {
	queueMap := t.buildQueueHandlers(ctx)
	for queueName, handler := range queueMap {
		q := queueName
		h := handler
		go func() {
			fmt.Println("[Worker] Listening on:", q)
			err := t.mqClient.Consume(q, func(msg amqp.Delivery) {
				if err := h(msg); err != nil {
					fmt.Println("[Worker] Error handling", q, ":", err)
				}
			})
			if err != nil {
				fmt.Println("[Worker] Failed to consume from", q, ":", err)
			}
		}()
	}

	fmt.Println("[Worker] All queues are listening. Press CTRL+C to stop.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("[Worker] Shutdown gracefully.")
	return nil
}

func (t *TaskWorker) buildQueueHandlers(ctx context.Context) map[string]func(amqp.Delivery) error {
	return map[string]func(amqp.Delivery) error{
		t.conf.ArticleQueue: handleStoreArticle(ctx, t.esClient, t.conf.ArticleQueue),
	}
}


func handleStoreArticle(ctx context.Context, es *elasticsearch.Elasticsearch, indexName string) func(amqp.Delivery) error {
	return func(msg amqp.Delivery) error {
		var article structs.PayloadMessageArticle
		if err := json.Unmarshal(msg.Body, &article); err != nil {
			return err
		}
		newArticle := article.NewArticle()
		if err := es.StoreToElasticsearch(ctx, newArticle); err != nil {
			return  err
		}
		return nil
	}
}




func NewTaskWorker(conf *config.Config, mq *rabbitmq.RabbitMQClient, es *elasticsearch.Elasticsearch) *TaskWorker {
	return &TaskWorker{
		conf:     conf,
		mqClient: mq,
		esClient: es,
	}
}

