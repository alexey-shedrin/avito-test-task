package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	HttpServer       HttpServer       `yaml:"http_server"`
	Database         Database         `yaml:"database"`
	GrpcServer       GrpcServer       `yaml:"grpc_server"`
	PrometheusServer PrometheusServer `yaml:"prometheus_server"`
}

type HttpServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type GrpcServer struct {
	Port string `yaml:"port"`
}

type PrometheusServer struct {
	Port string `yaml:"port"`
}

func New() *Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./configs/config.yaml"
	}

	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return cfg
}
