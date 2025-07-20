package main

import (
	"log"
	"os"

	"github.com/dika22/news-service/internal/domain/article/repository"
	"github.com/dika22/news-service/internal/domain/article/usecase"
	authorRepository "github.com/dika22/news-service/internal/domain/author/repository"
	"github.com/dika22/news-service/package/config"
	"github.com/dika22/news-service/package/connection/cache"
	database "github.com/dika22/news-service/package/connection/database/postgre"
	"github.com/dika22/news-service/package/connection/elasticsearch"
	rabbitmq "github.com/dika22/news-service/package/rabbit-mq"

	api "github.com/dika22/news-service/cmd/api"
	"github.com/dika22/news-service/cmd/worker"

	"github.com/urfave/cli/v2"

	cacheRepo "github.com/dika22/news-service/internal/domain/article/repository/cache"

	"github.com/dika22/news-service/package/validator"
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
  validate := validator.NewValidator()

  usecase := usecase.NewsUsecase(articleRepo, authorRepo, mqClient, esClient, conf, cacheArticleRepo)
  cmds := []*cli.Command{}
  cmds = append(cmds, api.ServeAPI(conf, validate, usecase)...)
  // cmds = append(cmds, migrate.NewMigrate(dbConn)...)
  cmds = append(cmds, worker.StartWorker(conf, mqClient, esClient, cacheArticleRepo)...)
  app := &cli.App{
    Name: "news-service",
    Commands: cmds,
  }

  if err := app.Run(os.Args); err != nil {
    panic(err)
  }
}
