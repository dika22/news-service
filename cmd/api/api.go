package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dika22/news-service/cmd/middleware"
	"github.com/dika22/news-service/internal/domain/article/delivery"
	"github.com/dika22/news-service/internal/domain/article/usecase"
	"github.com/dika22/news-service/package/config"
	"github.com/dika22/news-service/package/logger"

	"os/signal"

	"github.com/dika22/news-service/metrics"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/urfave/cli/v2"

	_ "github.com/dika22/news-service/docs"

	"github.com/dika22/news-service/package/validator"
	echoMiddlerware "github.com/labstack/echo/v4/middleware"
)

const CmdServeHTTP = "serve-http"

type HTTP struct{
	usecase usecase.IArticle
	conf *config.Config
	v *validator.Validator
}

func (h HTTP) ServeAPI(c *cli.Context) error  {
	if err := logger.SetLogger(); err != nil {
		log.Printf("error logger %v", err)
	}
	// Register metrics
	metrics.Register()
	e := echo.New();
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health-check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok!")
	})

	e.Use(echoMiddlerware.CORSWithConfig(echoMiddlerware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
	
	// Configurable rate limiter
	ipLimiter := middleware.RateLimiterMiddleware(
		cast.ToInt(h.conf.RateLimitMaxRequest), 
		time.Duration(cast.ToInt(h.conf.RateLimitInterval)) * time.Second, 
		cast.ToFloat64(h.conf.RateLimitJitter),
	) 
	e.Use(ipLimiter.Middleware())

	articleAPI := e.Group("api/v1/articles")
	articleAPI.Use(middleware.LoggerMiddleware)
	articleAPI.Use(middleware.MonitoringMiddleware)

	delivery.NewArticleHTTP(articleAPI, h.usecase, h.v)

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

func ServeAPI(conf *config.Config, v *validator.Validator, usecase usecase.IArticle) []*cli.Command {
	h := &HTTP{conf: conf, usecase: usecase, v: v}
	return []*cli.Command{
		{
			Name: CmdServeHTTP,
			Usage: "Serve News Service",
			Action: h.ServeAPI,
		},
	}
}