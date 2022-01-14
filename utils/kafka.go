package utils

import (
	"Dp-218_GO_micro/configs"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"strconv"
	"sync"
	"time"
)

const TopicName = "important"
const ClientID = "some_client"
const GroupConsumer = "some_group"

var kafkaVersion = sarama.V3_0_0_0

func CheckKafka() {

	broker, err := connectToBroker(configs.KAFKA_BROKER)
	if err != nil {
		log.Fatalln("Failed to connect to kafka broker:", err)
		return
	}

	err = createTopic([]string{configs.KAFKA_BROKER}, TopicName, 1, 1)
	if err != nil && err.(*sarama.TopicError).Err != sarama.ErrTopicAlreadyExists {
		log.Fatalln("Failed to create kafka topic:", err)
		return
	}

	producer := createProducer([]string{configs.KAFKA_BROKER}, ClientID)
	for i := 0; i < 10; i++ {
		_ = sendMessage(producer, TopicName, "Hello there"+strconv.Itoa(i))
	}

	group := createConsumerGroup([]string{configs.KAFKA_BROKER}, ClientID, GroupConsumer)
	ctx, cancel := context.WithCancel(context.Background())
	consumeMessages(ctx, group, TopicName)
	cancel()

	group.Close()
	producer.Close()
	broker.Close()

}

func connectToBroker(brokerAddr string) (*sarama.Broker, error) {
	config := sarama.NewConfig()
	config.Version = kafkaVersion

	broker := sarama.NewBroker(brokerAddr)

	connAttempts := 5
	connTimeout := 2 * time.Second
	var isConnected bool
	var err error
	for connAttempts > 0 {
		_ = broker.Open(config)
		isConnected, err = broker.Connected()
		if isConnected && err == nil {
			break
		}

		log.Printf("Kafka. Trying to connect, attempts left: %d", connAttempts)
		time.Sleep(connTimeout)
		connAttempts--
	}

	if !isConnected || err != nil {
		return broker, fmt.Errorf("kafka failed to connect")
	}

	return broker, err
}

func createTopic(brokerList []string, topicName string, nPartitions int32, replicas int16) error {
	config := sarama.NewConfig()
	config.Version = kafkaVersion

	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		return err
	}
	defer func() { _ = admin.Close() }()

	err = admin.CreateTopic(topicName, &sarama.TopicDetail{
		NumPartitions:     nPartitions,
		ReplicationFactor: replicas,
	}, false)

	return err
}

func createProducer(brokerList []string, clientID string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.ClientID = clientID

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func sendMessage(producer sarama.SyncProducer, topic, message string) error {
	_, _, err := producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	})

	return err
}

func createConsumerGroup(brokerList []string, clientID string, groupName string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = kafkaVersion
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.ClientID = clientID

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupName, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v\n", err)
	}
	return consumerGroup
}

func consumeMessages(ctx context.Context, group sarama.ConsumerGroup, topic string) {
	consumer := Consumer{
		ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := group.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Fatalf("Error from consumer: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}

			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	wg.Wait()
}

type Consumer struct {
	ready chan bool
}

func (consumer *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

func (consumer *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Kafka: value=%s, time=%v, topic=%s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
