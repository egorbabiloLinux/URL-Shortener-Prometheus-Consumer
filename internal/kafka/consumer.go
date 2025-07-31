package kafka

import (
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	log *slog.Logger
	*kafka.Consumer	
}

type ConsumerConfig interface {
	Get(key string) (string, bool)
}

func NewConsumer(cfg ConsumerConfig, log *slog.Logger, topics ...string) (*KafkaConsumer, error) {
	get := func(key string) string {
		val, _ := cfg.Get(key)
		return val
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
	"bootstrap.servers":                    get("bootstrap.servers"),
	"group.id":                             get("group.id"), 
	"auto.offset.reset":                    get("auto.offset.reset"), 

	"security.protocol":                    "SASL_SSL",
	"sasl.mechanisms":                      "PLAIN",
	"sasl.username":                        get("sasl.username"),
	"sasl.password":                        get("sasl.password"),

	"ssl.keystore.location":                get("ssl.keystore.location"),
	"ssl.keystore.password":                get("ssl.keystore.password"),
	"ssl.truststore.location":              get("ssl.truststore.location"),
	"ssl.truststore.password":              get("ssl.truststore.password"),
	"ssl.endpoint.identification.algorithm": get("ssl.endpoint.identification.algorithm"),
	})
	if err != nil {
		return nil, err
	}
	defer c.Close()

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		log: log,
		Consumer: c,
	}, nil
}