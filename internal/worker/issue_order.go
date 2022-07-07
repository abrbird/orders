package worker

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	cnfg "github.com/abrbird/orders/config"
	"github.com/abrbird/orders/internal/broker/kafka"
	"github.com/abrbird/orders/internal/metrics"
	"github.com/abrbird/orders/internal/models"
	rpstr "github.com/abrbird/orders/internal/repository"
	srvc "github.com/abrbird/orders/internal/service"
	"github.com/pkg/errors"
	"log"
)

type MarkOrderIssuedHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
	metrics    metrics.Metrics
	config     *cnfg.Config
}

func (i *MarkOrderIssuedHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *MarkOrderIssuedHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (i *MarkOrderIssuedHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := context.Background()

		if msg.Topic != i.config.Kafka.IssueOrderTopics.MarkOrderIssued {
			log.Printf(
				"topic names does not match: expected - %s, got %s\n",
				i.config.Kafka.IssueOrderTopics.MarkOrderIssued,
				msg.Topic,
			)
			continue
		}

		var issueOrderMessage kafka.IssueOrderMessage
		err := json.Unmarshal(msg.Value, &issueOrderMessage)
		if err != nil {
			i.metrics.Error()
			log.Print("Unmarshall failed: value=%v, err=%v", string(msg.Value), err)
			continue
		}

		log.Printf("consumer %s: <- %s: %v",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.MarkOrderIssued,
			issueOrderMessage,
		)

		err = i.service.Order().MarkOrderIssued(
			ctx,
			i.repository.Order(),
			issueOrderMessage.Order.Id,
		)
		if err != nil {
			i.metrics.Error()
			if errors.Is(err, models.RetryError) {
				err = i.RetryMarkOrderIssued(issueOrderMessage)
				if err != nil {
					i.metrics.KafkaError()
					log.Println(err)
				} else {
					log.Printf(
						"consumer %s: -> %s: %v",
						i.config.Application.Name,
						i.config.Kafka.IssueOrderTopics.MarkOrderIssued,
						issueOrderMessage,
					)
				}
			} else {
				i.metrics.KafkaError()
				log.Println(err)
			}
			continue
		}

		err = i.SendConfirmIssueOrder(issueOrderMessage)
		if err != nil {
			i.metrics.KafkaError()
			log.Println(err)
		} else {
			log.Printf(
				"consumer %s: -> %s: %v",
				i.config.Application.Name,
				i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder,
				issueOrderMessage,
			)
		}
	}
	return nil
}

func (i *MarkOrderIssuedHandler) RetryMarkOrderIssued(message kafka.IssueOrderMessage) error {
	message.Base.SenderServiceName = i.config.Application.Name
	message.Base.Attempt += 1

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.MarkOrderIssued, message)
	if err != nil {
		return models.BrokerSendError(err)
	}

	if kerr != nil {
		return models.BrokerSendError(err)
	}
	_ = part
	_ = offs

	return nil
}

func (i *MarkOrderIssuedHandler) SendConfirmIssueOrder(message kafka.IssueOrderMessage) error {
	message.Base.SenderServiceName = i.config.Application.Name

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder, message)
	if err != nil {
		return models.BrokerSendError(err)
	}

	if kerr != nil {
		return models.BrokerSendError(err)
	}
	_ = part
	_ = offs

	return nil
}
