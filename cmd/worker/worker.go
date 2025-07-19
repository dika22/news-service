package worker

import (
	"context"
	"news-service/cmd/worker/task"
	"news-service/package/config"
	"news-service/package/connection/elasticsearch"
	rabbitmq "news-service/package/rabbit-mq"

	"github.com/urfave/cli/v2"
)
const CmdStartWorker = "start-worker"

type Worker struct {
	conf     *config.Config
	mqClient *rabbitmq.RabbitMQClient
	esClient *elasticsearch.Elasticsearch

}

func (w *Worker) StartWorker(c *cli.Context) error {
	ctx := context.Background()
	worker := task.NewTaskWorker(w.conf, w.mqClient, w.esClient)
	return worker.Run(ctx)
}


func StartWorker(conf *config.Config, mqClient *rabbitmq.RabbitMQClient, esClient *elasticsearch.Elasticsearch) []*cli.Command {
	w := &Worker{
		conf:     conf,
		mqClient: mqClient,
		esClient: esClient,
	}
	return []*cli.Command{
		{
			Name:  CmdStartWorker,
			Usage: "Start background worker to consume tasks from queue",
			Action: w.StartWorker,
		},
	}
}
