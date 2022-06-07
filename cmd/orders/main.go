package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.ozon.dev/zBlur/homework-3/orders/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/db"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/repository/sql_repository"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/service/implemented_service"
	wrkr "gitlab.ozon.dev/zBlur/homework-3/orders/internal/worker"
)

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

	repository := sql_repository.New(dbConnPool)
	service := implemented_service.New()

	//tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	//defer closer.Close()
	//opentracing.SetGlobalTracer(
	//	tracer,
	//)

	fmt.Println("Start orders...")
	fmt.Println("config", cfg)
	fmt.Println("repository", repository)
	fmt.Println("service", service)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	worker, err := wrkr.New(cfg, repository, service)
	err = worker.StartConsuming(ctx)

	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
