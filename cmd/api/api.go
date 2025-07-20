package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"news-service/cmd/middleware"
	"news-service/internal/domain/article/delivery"
	"news-service/internal/domain/article/usecase"
	"news-service/package/config"
	"news-service/package/logger"
	"os"
	"time"

	"os/signal"

	"news-service/metrics"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/urfave/cli/v2"

	_ "news-service/docs"

	echoMiddlerware "github.com/labstack/echo/v4/middleware"
)

const CmdServeHTTP = "serve-http"

type HTTP struct{
	usecase usecase.IArticle
	conf *config.Config
}

func (h HTTP) ServeAPI(c *cli.Context) error  {
	if err := logger.SetLogger(); err != nil {
		log.Printf("error logger %v", err)
	}
	// Register metrics
	metrics.Register()
	e := echo.New();
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "pong")
	})

	e.Use(echoMiddlerware.CORSWithConfig(echoMiddlerware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	
	fmt.Println("Rate limit config:", h.conf.RateLimitMaxToken, h.conf.RateLimitInterval, h.conf.RateLimitJitter)

	// Configurable rate limiter
	// 5 req / sec with 20% jitter
	rl := middleware.NewRateLimiter(
		cast.ToInt(h.conf.RateLimitMaxToken), 
		time.Duration(cast.ToInt(h.conf.RateLimitInterval)) * time.Second, 
		cast.ToFloat64(h.conf.RateLimitJitter),
	) 

	e.Use(middleware.RateLimiterMiddleware(rl))
	articleAPI := e.Group("api/v1/articles")
	articleAPI.Use(middleware.LoggerMiddleware)
	articleAPI.Use(middleware.MonitoringMiddleware)

	delivery.NewArticleHTTP(articleAPI, h.usecase)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(fmt.Sprintf(":%v", h.conf.AppPort)); err != nil {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	return nil	
}

func ServeAPI(conf *config.Config, usecase usecase.IArticle) []*cli.Command {
	h := &HTTP{conf: conf, usecase: usecase}
	return []*cli.Command{
		{
			Name: CmdServeHTTP,
			Usage: "Serve Document Service",
			Action: h.ServeAPI,
		},
	}
}