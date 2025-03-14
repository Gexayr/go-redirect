package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	MySQL    MySQLConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	Port string
}

type MySQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Queue    string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("mysql.host", "mysql")
	viper.SetDefault("mysql.port", "3306")
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "root")
	viper.SetDefault("mysql.dbname", "platform_db")
	viper.SetDefault("rabbitmq.host", "rabbitmq")
	viper.SetDefault("rabbitmq.port", "5672")
	viper.SetDefault("rabbitmq.user", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.queue", "request_queue")

	// Read environment variables
	viper.BindEnv("mysql.host", "MYSQL_HOST")
	viper.BindEnv("mysql.port", "MYSQL_PORT")
	viper.BindEnv("mysql.user", "MYSQL_USER")
	viper.BindEnv("mysql.password", "MYSQL_PASSWORD")
	viper.BindEnv("mysql.dbname", "MYSQL_DBNAME")
	viper.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq.user", "RABBITMQ_USER")
	viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")
	viper.BindEnv("rabbitmq.queue", "RABBITMQ_QUEUE")

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
} 