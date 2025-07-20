package task

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dika22/news-service/package/config"
	"github.com/dika22/news-service/package/connection/elasticsearch"
	rabbitmq "github.com/dika22/news-service/package/rabbit-mq"
	"github.com/dika22/news-service/package/structs"

	cacheRepo "github.com/dika22/news-service/internal/domain/article/repository/cache"
	"github.com/streadway/amqp"
)

type TaskWorker struct{
	conf     *config.Config
	mqClient *rabbitmq.RabbitMQClient
	esClient *elasticsearch.Elasticsearch
	cache    cacheRepo.CacheRepository
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
		t.conf.ArticleQueue: t.handleStoreArticle(ctx, t.esClient, t.conf.ArticleQueue),
	}
}


func (t *TaskWorker) handleStoreArticle(ctx context.Context, es *elasticsearch.Elasticsearch, indexName string) func(amqp.Delivery) error {
	return func(msg amqp.Delivery) error {
		var article structs.PayloadMessageArticle
		if err := json.Unmarshal(msg.Body, &article); err != nil {
			return err
		}
		newArticle := article.NewArticle()
		if err := es.StoreToElasticsearch(ctx, newArticle); err != nil {
			return  err
		}

		if err := t.cache.DeleteArticleKeys(ctx); err != nil {
			fmt.Println("error logs cache delete article keys", err)
		}
		return nil
	}
}


func NewTaskWorker(conf *config.Config, mq *rabbitmq.RabbitMQClient, 
	es *elasticsearch.Elasticsearch,
	cache cacheRepo.CacheRepository) *TaskWorker {
	return &TaskWorker{
		conf:     conf,
		mqClient: mq,
		esClient: es,
		cache:    cache,
	}
}

