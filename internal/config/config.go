package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	BootstrapServers 			 string `mapstructure:"bootstrap_servers" validate:"required"`
	GroupId 					 string `mapstructure:"group_id" validate:"required"`
	AutoOffsetReset 			 string `mapstructure:"auto_offset_reset" validate:"required"`
	SASLUsername 				 string `mapstructure:"sasl_username" validate:"required"`
	SASLPassword 				 string `mapstructure:"sasl_password" validate:"required"`
	SSLKeystoreLocation 		 string `mapstructure:"ssl_keystore_location" validate:"required"`
	SSLKeystorePassword 		 string `mapstructure:"ssl_keystore_password" validate:"required"`
	SSLTruststoreLocation 		 string `mapstructure:"ssl_truststore_location" validate:"required"`
	SSLTruststorePassword 		 string `mapstructure:"ssl_truststore_password" validate:"required"`
	SSLEndpointIdentificationAlg string `mapstructure:"ssl_endpoint_identification_algorithm" validate:"required"`
}

func (c Config) Get(key string) (string, bool) {
	switch key {
	case "bootstrap.servers":
		return c.BootstrapServers, true
	case "group.id":
		return c.GroupId, true 
	case "auto.offset.reset":
		return c.AutoOffsetReset, true
	case "sasl.username":
		return c.SASLUsername, true
	case "sasl.password":
		return c.SASLPassword, true
	case "ssl.keystore.location":
		return c.SSLKeystoreLocation, true
	case "ssl.keystore.password":
		return c.SSLKeystorePassword, true
	case "ssl.truststore.location":
		return c.SSLTruststoreLocation, true
	case "ssl.truststore.password":
		return c.SSLTruststorePassword, true
	case "ssl.endpoint.identification.algorithm":
		return c.SSLEndpointIdentificationAlg, true
	default:
		return "", false
	}
}

func MustLoad() *Config {
	v := viper.New()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		panic("config file does not exists: " + configPath)
	}

	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	v.AutomaticEnv()

	v.SetDefault("env", "local")

	envBindings := map[string]string{
		"group_id": 							 "CONSUMER_GROUP_ID",
		"auto_offset_reset": 					 "KAFKA_AUTO_OFFSET_RESET",
		"kafka_producer.bootstrap_servers":      "KAFKA_BOOTSTRAP_SERVERS",
		"kafka_producer.sasl_username":          "KAFKA_SASL_USERNAME",
		"kafka_producer.sasl_password":          "KAFKA_SASL_PASSWORD",
		"kafka_producer.ssl_keystore_location":  "KAFKA_SSL_KEYSTORE_LOCATION",
		"kafka_producer.ssl_keystore_password":  "KAFKA_SSL_KEYSTORE_PASSWORD",
		"kafka_producer.ssl_truststore_location":"KAFKA_SSL_TRUSTSTORE_LOCATION",
		"kafka_producer.ssl_truststore_password":"KAFKA_SSL_TRUSTSTORE_PASSWORD",
	}

	for key, envVar := range envBindings {
		if val := os.Getenv(envVar); val != "" {
			v.Set(key, val)
		}
	}

	var cfg Config 

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config: %s", err)
	}

	if err := validator.New().Struct(cfg); err != nil {
		log.Fatalf("error validating config: %s", err)
	}
	
	return &cfg
}