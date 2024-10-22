package http

import (
	"context"
	"fmt"
	"generic-integration-platform/internal/infra/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func newAppServer() *gin.Engine {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s - %s \"%s %s\" %d %s \"%s\" %s\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Request.Method,
			param.Request.URL.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	return r
}

func NewServer(lc fx.Lifecycle, config *config.Config) *gin.Engine {
	r := newAppServer()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				srv.ListenAndServe()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Close()
			return nil
		},
	})

	return r
}
