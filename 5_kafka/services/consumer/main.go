package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

var (
	kafkaHost string
	topic     string
)

func init() {
	kafkaHost = os.Getenv("KAFKA_HOST")
	topic = os.Getenv("TOPIC")
}

func main() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaHost, topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal("failed to close connection:", err)
		}
	}()

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	for {
		msg, err := conn.ReadMessage(1e6) // 1MB max
		if err != nil {
			break
		}
		fmt.Printf("message received. topiс: %s, partition: %d, key: %s, value: %s\n", msg.Topic, msg.Partition, string(msg.Key), string(msg.Value))
	}
}
