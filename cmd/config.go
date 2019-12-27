package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	Etcd     []etcdConfig `yaml:"etcd"`
	HTTPPort string       `yaml:"httpport"`
}

type etcdConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// ReadConfig reads config file based on env
func ReadConfig(path string) config {
	viper.AddConfigPath(path)
	fmt.Println(os.Getenv("CT_DNS_ENV"))
	if os.Getenv("CT_DNS_ENV") == "PRODUCTION" {
		viper.SetConfigName("production")
	} else {
		viper.SetConfigName("development")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Failed to read config", err)
	}
	var c config
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalln("Failed to unmarshal config", err)
	}
	return c
}
