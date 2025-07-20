package worker

import (
	"context"

	"github.com/dika22/news-service/cmd/worker/task"
	"github.com/dika22/news-service/package/config"
	"github.com/dika22/news-service/package/connection/elasticsearch"
	rabbitmq "github.com/dika22/news-service/package/rabbit-mq"

	cacheRepo "github.com/dika22/news-service/internal/domain/article/repository/cache"
	"github.com/urfave/cli/v2"
)
const CmdStartWorker = "start-worker"

type Worker struct {
	conf     *config.Config
	mqClient *rabbitmq.RabbitMQClient
	esClient *elasticsearch.Elasticsearch
	cache    cacheRepo.CacheRepository
}

func (w *Worker) StartWorker(c *cli.Context) error {
	ctx := context.Background()
	worker := task.NewTaskWorker(w.conf, w.mqClient, w.esClient, w.cache)
	return worker.Run(ctx)
}


func StartWorker(conf *config.Config, mqClient *rabbitmq.RabbitMQClient, 
	esClient *elasticsearch.Elasticsearch, 
	cache cacheRepo.CacheRepository) []*cli.Command {
	w := &Worker{
		conf:     conf,
		mqClient: mqClient,
		esClient: esClient,
		cache:    cache,
	}
	return []*cli.Command{
		{
			Name:  CmdStartWorker,
			Usage: "Start background worker to consume tasks from queue",
			Action: w.StartWorker,
		},
	}
}
