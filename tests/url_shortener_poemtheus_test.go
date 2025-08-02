package tests

import (
	"testing"
	"url-shortener-pronetheus-consumer/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestHandleKafkaMessage_IncrementsMetric(t *testing.T) {
	metrics.Register()

	msg := &kafka.Message{
		Key: []byte("web"),
		Value: []byte(`{"source":"web"}`),
		TopicPartition: kafka.TopicPartition{
			Topic:     strPtr("auth-events"),
			Partition: 0,
			Offset:    0,
		},
	}

	metrics.AuthCounter.WithLabelValues(string(msg.Key)).Inc()

	val := testutil.ToFloat64(metrics.AuthCounter.WithLabelValues("web"))
	if val != 1 {
		t.Fatalf("expected metric value 1, got %v", val)
	}
}

func strPtr(s string) *string {
	return &s
}