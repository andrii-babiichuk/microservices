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
	topic = os.Getenv("POSTGRES_PASSWORD")
}

func main() {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaHost, topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
	}()
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	for i := 0; ; i++ {
		_, err = conn.WriteMessages(
			kafka.Message{Key: []byte(fmt.Sprintf("%d", i)), Value: []byte(fmt.Sprintf("message %d", i))},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
		time.Sleep(time.Second * 30)
	}
}
