package config

import (
	"os"

	"github.com/guanw/ct-dns/pkg/logging"
	"github.com/spf13/viper"
)

// Config contains config for ct-dns service
type Config struct {
	Etcd     []etcdConfig `yaml:"etcd"`
	HTTPPort string       `yaml:"httpport"`
	GRPCPort string       `yaml:"grpcport"`
}

type etcdConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// ReadConfig reads config file based on env
func ReadConfig(path string) Config {
	viper.AddConfigPath(path)
	logging.GetLogger().WithField("environment", os.Getenv("CT_DNS_ENV")).Info("Initializing config...")
	if os.Getenv("CT_DNS_ENV") == "PRODUCTION" {
		viper.SetConfigName("production")
	} else {
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
