package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Postgres PostgresConfig
	Port     string
	Redis    RedisConfig
	Kafka    KafkaConfig
}

type KafkaConfig struct {
	GroupID   string
	Url       string
	TopicName string
}

type PostgresConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Db       int
	Password string
}

func SLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisDb := os.Getenv("REDIS_DB")
	kafkaUrl := os.Getenv("KAFKA_URL")
	kafkaTopic := os.Getenv("KAFKA_TOPIC_NAME")
	kafkaGroup := os.Getenv("KAFKA_GROUP_ID")

	rDb, err := strconv.Atoi(redisDb)
	if err != nil {
		log.Fatal("Error parsing redis db")
	}

	port := os.Getenv("PORT")

	cfg := Config{
		Postgres: PostgresConfig{
			Host:     dbHost,
			Port:     dbPort,
			DbName:   dbName,
			User:     dbUser,
			Password: dbPassword,
		},
		Port: port,
		Redis: RedisConfig{
			Host:     redisHost,
			Port:     redisPort,
			Db:       rDb,
			Password: redisPassword,
		},
		Kafka: KafkaConfig{
			GroupID:   kafkaGroup,
			Url:       kafkaUrl,
			TopicName: kafkaTopic,
		},
	}

	return &cfg

}
