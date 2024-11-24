package config

import (
	"os"
)

const (
	DEVELOPER    = "developer"
	HOMOLOGATION = "homologation"
	PRODUCTION   = "production"
)

type Config struct {
	PORT          string `json:"port"`
	Mode          string `json:"mode"`
	MongoDBConfig `json:"mongo_config"`
	*GroqConfig
	*RMQConfig
}
type MongoDBConfig struct {
	MDB_URI                string `json:"mdb_uri"`
	MDB_NAME               string `json:"mdb_name"`
	MDB_DEFAULT_COLLECTION string `json:"mdb_default_collection"`
}

type AsaasConfig struct {
	URL_ASAAS       string `json:"url_asaas"`
	ASAAS_API_KEY   string `json:"asas_api_key"`
	ASAAS_WALLET_ID string `json:"asas_wallet_id"`
	ASAAS_TIMEOUT   int    `json:"asas_timeout"`
}

type GroqConfig struct {
	API_KEY string `json:"api_key"`
	URL     string `json:"url"`
}

type RMQConfig struct {
	RMQ_URI                  string `json:"rmq_uri"`
	RMQ_MAXX_RECONNECT_TIMES int    `json:"rmq_maxx_reconnect_times"`
}

type ConsumerConfig struct {
	ExchangeName  string `json:"exchange_name"`
	ExchangeType  string `json:"exchange_type"`
	RoutingKey    string `json:"routing_key"`
	QueueName     string `json:"queue_name"`
	ConsumerName  string `json:"consumer_name"`
	ConsumerCount int    `json:"consumer_count"`
	PrefetchCount int    `json:"prefetch_count"`
	Reconnect     struct {
		MaxAttempt int `json:"max_attempt"`
		Interval   int `json:"interval"`
	}
}

func NewConfig() *Config {
	conf := defaultConf()

	if port := os.Getenv("SRV_PORT"); port != "" {
		conf.PORT = port
	}

	if mode := os.Getenv("SRV_MODE"); mode != "" {
		conf.Mode = mode
	}

	if uri := os.Getenv("SRV_MDB_URI"); uri != "" {
		conf.MDB_URI = uri
	}

	if name := os.Getenv("SRV_MDB_NAME"); name != "" {
		conf.MDB_NAME = name
	}

	if collection := os.Getenv("SRV_MDB_DEFAULT_COLLECTION"); collection != "" {
		conf.MDB_DEFAULT_COLLECTION = collection
	}

	if groqApiKey := os.Getenv("GROQ_API_KEY"); groqApiKey != "" {
		conf.GroqConfig.API_KEY = groqApiKey
	}

	if groqUrl := os.Getenv("GROQ_URL"); groqUrl != "" {
		conf.GroqConfig.URL = groqUrl
	}

	SRV_RMQ_URI := os.Getenv("SRV_RMQ_URI")
	if SRV_RMQ_URI != "" {
		conf.RMQConfig.RMQ_URI = SRV_RMQ_URI

	}

	return conf
}

func defaultConf() *Config {
	defaultConf := &Config{
		PORT: "8080",
		Mode: DEVELOPER,

		MongoDBConfig: MongoDBConfig{
			MDB_DEFAULT_COLLECTION: "cfSera",
		},

		GroqConfig: &GroqConfig{
			URL: "https://api.groq.com/openai/v1/chat/completions",
		},

		RMQConfig: &RMQConfig{

			RMQ_MAXX_RECONNECT_TIMES: 3,
		},
	}

	return defaultConf
}
