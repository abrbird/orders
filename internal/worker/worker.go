package worker

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/abrbird/orders/internal/metrics"
	"log"
	"time"

	cnfg "github.com/abrbird/orders/config"
	"github.com/abrbird/orders/internal/broker/kafka"
	rpstr "github.com/abrbird/orders/internal/repository"
	srvc "github.com/abrbird/orders/internal/service"
)

type OrdersTrackingWorker struct {
	config          *cnfg.Config
	repository      rpstr.Repository
	service         srvc.Service
	metrics         metrics.Metrics
	producer        sarama.SyncProducer
	markOrderIssued *MarkOrderIssuedHandler
}

func New(cfg *cnfg.Config, repository rpstr.Repository, service srvc.Service, metrics metrics.Metrics) (*OrdersTrackingWorker, error) {

	brokerConfig := kafka.NewConfig()
	producer, err := kafka.NewSyncProducer(cfg.Kafka.Brokers.String(), brokerConfig)
	if err != nil {
		return nil, err
	}

	worker := &OrdersTrackingWorker{
		config:     cfg,
		repository: repository,
		service:    service,
		metrics:    metrics,
		producer:   producer,
		markOrderIssued: &MarkOrderIssuedHandler{
			producer:   producer,
			repository: repository,
			service:    service,
			metrics:    metrics,
			config:     cfg,
		},
	}

	return worker, nil
}

func (w *OrdersTrackingWorker) StartConsuming(ctx context.Context) error {

	brokerConfig := kafka.NewConfig()
	markOrderIssued, err := sarama.NewConsumerGroup(
		w.config.Kafka.Brokers.String(),
		fmt.Sprintf("%s%sCG", w.config.Application.Name, w.config.Kafka.IssueOrderTopics.MarkOrderIssued),
		brokerConfig,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			err := markOrderIssued.Consume(ctx, []string{w.config.Kafka.IssueOrderTopics.MarkOrderIssued}, w.markOrderIssued)
			if err != nil {
				log.Printf("%s consumer error: %v", w.config.Kafka.IssueOrderTopics.MarkOrderIssued, err)
				time.Sleep(time.Second * 5)
			}
		}
	}()
	go func() {
		for err := range markOrderIssued.Errors() {
			log.Println(err)
		}
	}()

	return nil
}
