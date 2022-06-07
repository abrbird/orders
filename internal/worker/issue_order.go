package worker

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	cnfg "gitlab.ozon.dev/zBlur/homework-3/orders/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/broker/kafka"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/models"
	rpstr "gitlab.ozon.dev/zBlur/homework-3/orders/internal/repository"
	srvc "gitlab.ozon.dev/zBlur/homework-3/orders/internal/service"
	"log"
)

type MarkOrderIssuedHandler struct {
	producer   sarama.SyncProducer
	repository rpstr.Repository
	service    srvc.Service
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

		log.Printf("consumer %s: -> %s: %v",
			i.config.Application.Name,
			i.config.Kafka.IssueOrderTopics.MarkOrderIssued,
			msg.Value,
		)

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
			log.Print("Unmarshall failed: value=%v, err=%v", string(msg.Value), err)
			continue
		}

		ctx := context.Background()
		orderRetrieved := i.service.Order().Retrieve(
			ctx,
			i.repository.Order(),
			issueOrderMessage.Order.Id,
		)

		if orderRetrieved.Error != nil {
			log.Printf("error on message processing: %v", err)
			i.RetryMarkOrderIssued(issueOrderMessage)
			continue
		}
		orderRetrieved.Order.Status = models.Issued

		err = i.service.Order().Update(
			ctx,
			i.repository.Order(),
			orderRetrieved.Order,
		)
		if err != nil {
			log.Printf("error on order update: %v", err)
			i.RetryMarkOrderIssued(issueOrderMessage)
			continue
		}

		i.SendConfirmIssueOrder(issueOrderMessage)
	}
	return nil
}

func (i *MarkOrderIssuedHandler) RetryMarkOrderIssued(message kafka.IssueOrderMessage) {
	message.Base.SenderServiceName = i.config.Application.Name
	message.Base.Attempt += 1

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.IssueOrder, message)
	if err != nil {
		log.Printf("can not send message: %v", err)
		return
	}

	if kerr != nil {
		log.Printf("can not send message: %v", kerr)
		return
	}

	log.Printf("consumer %s: %v -> %v", i.config.Application.Name, part, offs)
	return
}

func (i *MarkOrderIssuedHandler) SendConfirmIssueOrder(message kafka.IssueOrderMessage) {
	message.Base.SenderServiceName = i.config.Application.Name

	part, offs, kerr, err := kafka.SendMessage(i.producer, i.config.Kafka.IssueOrderTopics.ConfirmIssueOrder, message)
	if err != nil {
		log.Printf("can not send message: %v", err)
		return
	}

	if kerr != nil {
		log.Printf("can not send message: %v", kerr)
		return
	}

	log.Printf("consumer %s: %v -> %v", i.config.Application.Name, part, offs)
	return
}
