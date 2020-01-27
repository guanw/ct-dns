package config

import (
	"os"

	"github.com/guanw/ct-dns/pkg/logging"
	"github.com/spf13/viper"
)

// Config contains config for ct-dns service
type Config struct {
	Etcd     []EtcdConfig `yaml:"etcd"`
	HTTPPort string       `yaml:"httpport"`
	GRPCPort string       `yaml:"grpcport"`
	Redis    RedisConfig  `yaml:"redis"`
}

// EtcdConfig contains config for etcd cluster
type EtcdConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// RedisConfig contains config for redis cluster
type RedisConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// ReadConfig reads config file based on env
func ReadConfig(path string) Config {
	viper.AddConfigPath(path)
	logging.GetLogger().WithField("environment", os.Getenv("CT_DNS_ENV")).Info("Initializing config...")
	env := os.Getenv("CT_DNS_ENV")

	switch env {
	case "PRODUCTION":
		viper.SetConfigName("production")
	case "KUBERNETERS-REDIS":
		viper.SetConfigName("kubernetes-with-redis")
	default:
		viper.SetConfigName("development")
	}
	err := viper.ReadInConfig()
	if err != nil {
		logging.GetLogger().Fatalln("Failed to read config", err)
	}
	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		logging.GetLogger().Fatalln("Failed to unmarshal config", err)
	}
	return c
}
