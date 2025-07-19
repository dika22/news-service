package main

import (
	"log"
	"news-service/internal/domain/article/repository"
	"news-service/internal/domain/article/usecase"
	authorRepository "news-service/internal/domain/author/repository"
	"news-service/package/config"
	"news-service/package/connection/cache"
	database "news-service/package/connection/database/postgre"
	"news-service/package/connection/elasticsearch"
	rabbitmq "news-service/package/rabbit-mq"
	"os"

	api "news-service/cmd/api"
	"news-service/cmd/worker"

	"github.com/urfave/cli/v2"

	cacheRepo "news-service/internal/domain/article/repository/cache"
)

func main() {

  dbConf := config.NewDatabase()
  conf := config.NewConfig()
  dbConn := database.NewDatabase(dbConf)
  defer dbConn.Close()

  cacheConf := config.NewCache()
  cache := cache.NewRedis("web", cacheConf)

  mqClient, err := rabbitmq.NewRabbitMQClient(conf)
  if err != nil {
    log.Println("ERROR INIT RABBITMQ", err)
  }

  esClient, err := elasticsearch.NewElasticSearch(conf)
  if err != nil {
    log.Println("ERROR Connect Elasticsearch", err)
  }
  
  articleRepo := repository.NewsRepository(dbConn, cache)
  authorRepo := authorRepository.NewAuthorRepository(dbConn)
  cacheArticleRepo := cacheRepo.NewCacheRepository(cache)

  usecase := usecase.NewsUsecase(articleRepo, authorRepo, mqClient, esClient, conf, cacheArticleRepo)
  cmds := []*cli.Command{}
  cmds = append(cmds, api.ServeAPI(conf, usecase)...)
  // cmds = append(cmds, migrate.NewMigrate(dbConn)...)
  cmds = append(cmds, worker.StartWorker(conf, mqClient, esClient)...)
  app := &cli.App{
    Name: "news-service",
    Commands: cmds,
  }

  if err := app.Run(os.Args); err != nil {
    panic(err)
  }
}
