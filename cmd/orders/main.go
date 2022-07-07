package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/cache/redis_cache"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/metrics/prom_metrics"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"gitlab.ozon.dev/zBlur/homework-3/orders/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/db"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/repository/sql_repository"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/service/implemented_service"
	wrkr "gitlab.ozon.dev/zBlur/homework-3/orders/internal/worker"
)

type GoroutinesNum struct{}

func (g *GoroutinesNum) String() string {
	return strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
}

func main() {
	cfg, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	dbConnPool, err := db.New(cfg.Database, ctx)
	if err != nil {
		log.Fatal(err)
	}

	cache := redis_cache.New(cfg.Cache.Redis)

	repository := sql_repository.New(dbConnPool)
	service := implemented_service.New(cache)
	metrics := prom_metrics.New(cfg)

	//tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	//defer closer.Close()
	//opentracing.SetGlobalTracer(
	//	tracer,
	//)

	log.Printf("Start %s...", cfg.Application.Name)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker, err := wrkr.New(cfg, repository, service, metrics)
	if err != nil {
		log.Fatal(err)
	}
	err = worker.StartConsuming(ctx)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		g := &GoroutinesNum{}
		expvar.Publish("GoroutinesNum", g)

		err = http.ListenAndServe(
			fmt.Sprintf("%s:%v", cfg.Monitoring.Pprof.Host, cfg.Monitoring.Pprof.Port),
			nil,
		)
		if err != nil {
			log.Print(err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err = http.ListenAndServe(
			fmt.Sprintf("%s:%v", cfg.Monitoring.Metrics.Host, cfg.Monitoring.Metrics.Port),
			nil,
		)
		if err != nil {
			log.Print(err)
		}
	}()

	<-ctx.Done()
}
