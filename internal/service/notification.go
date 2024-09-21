package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

const timeoutFactor = 10

func NewNotificationService(topic, address string) *NotificationService {
	return &NotificationService{
		topic:   topic,
		address: address,
	}
}

type NotificationService struct {
	topic   string
	address string
}

func (s *NotificationService) Publish(message any) (err error) {
	partition := 0

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", s.address, s.topic, partition)
	if err != nil {
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(timeoutFactor * time.Second))
	if err != nil {
		return err
	}

	_, err = conn.Write(messageBytes)
	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}
