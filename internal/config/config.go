package config

import (
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	BootstrapServers 			 string   `mapstructure:"bootstrap_servers" validate:"required"`
	GroupId 					 string   `mapstructure:"group_id" validate:"required"`
	AutoOffsetReset 			 string   `mapstructure:"auto_offset_reset" validate:"required"`
	SASLUsername 				 string   `mapstructure:"sasl_username" validate:"required"`
	SASLPassword 				 string   `mapstructure:"sasl_password" validate:"required"`
	SSLKeyLocation 		 		 string   `mapstructure:"ssl_key_location" validate:"required"`
	SSLCertificateLocation 		 string   `mapstructure:"ssl_certificate_location" validate:"required"`
	SSLCaLocation		 		 string   `mapstructure:"ssl_ca_location" validate:"required"`
	SSLEndpointIdentificationAlg string   `mapstructure:"ssl_endpoint_identification_algorithm"`
	Topics 						 []string `mapstructure:"topics" validate:"required,dive,required"`
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
	case "ssl.key.location":
		return c.SSLKeyLocation, true
	case "ssl.certificate.location":
		return c.SSLCertificateLocation, true
	case "ssl.ca.location":
		return c.SSLCaLocation, true
	case "ssl.endpoint.identification.algorithm":
		if c.SSLEndpointIdentificationAlg == "" {
			return c.SSLEndpointIdentificationAlg, true
		}
		return "", false
	default:
		return "", false
	}
}

func MustLoad() *Config {
	err := godotenv.Load("./config/.env")
	if err != nil {
	 	log.Println(".env file not found or failed to load, skipping: " + err.Error())
	}

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
		"bootstrap_servers":           "KAFKA_BOOTSTRAP_SERVERS",
		"group_id":                    "CONSUMER_GROUP_ID",
		"auto_offset_reset":           "KAFKA_AUTO_OFFSET_RESET",
		"sasl_username":               "KAFKA_SASL_USERNAME",
		"sasl_password":               "KAFKA_SASL_PASSWORD",
		"ssl_key_location":       	   "KAFKA_SSL_KEY_LOCATION",
		"ssl_certificate_location":    "KAFKA_SSL_CERTIFICATE_LOCATION",
		"ssl_ca_location":     		   "KAFKA_SSL_CA_LOCATION",
		"topics":                      "TOPICS",
	}

	for key, envVar := range envBindings {
		if val := os.Getenv(envVar); val != "" {
			v.Set(key, val)
		}
	}

	if topicsEnv := os.Getenv(envBindings["topics"]); topicsEnv != "" {
		v.Set("topics", parseCSV(topicsEnv))
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

func parseCSV(str string) []string {
	parts := strings.Split(str, ",")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
