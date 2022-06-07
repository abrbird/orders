package main

import (
	"gitlab.ozon.dev/zBlur/homework-3/orders/config"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/broker/kafka"
	"log"
)

func main() {
	cfg, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	brokerCfg := kafka.NewConfig()
	syncProducer, err := kafka.NewSyncProducer(cfg.Kafka.Brokers.String(), brokerCfg)

	if err != nil {
		log.Fatal(err)
	}

	_ = syncProducer

	//for {
	//	d := models.Order{
	//		Id: time.Now().UnixNano(),
	//	}
	//	b, err := json.Marshal(d)
	//	if err != nil {
	//		log.Printf("wtf? %v", err)
	//		continue
	//	}
	//
	//	par, off, err := syncProducer.SendMessage(&sarama.ProducerMessage{
	//		Topic: cfg.Kafka.IssueOrderTopics.IssueOrder,
	//		Key:   sarama.StringEncoder(fmt.Sprintf("%v", d.Id)),
	//		Value: sarama.ByteEncoder(b),
	//	})
	//
	//	log.Printf("producer IssueOrder %v -> %v; %v", par, off, err)
	//	time.Sleep(time.Millisecond * 500)
	//
	//	if rand.Intn(10) == 9 {
	//		par, off, err = syncProducer.SendMessage(&sarama.ProducerMessage{
	//			Topic: cfg.Kafka.IssueOrderTopics.UndoIssueOrder,
	//			Key:   sarama.StringEncoder(fmt.Sprintf("%v", d.Id)),
	//			Value: sarama.ByteEncoder(b),
	//		})
	//		log.Printf("producer UndoIssueOrder %v -> %v; %v", par, off, err)
	//	}
	//}
}
