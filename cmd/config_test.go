package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadConfig(t *testing.T) {
	tests := []struct {
		env               string
		expectedHTTPPort  string
		expectedRedisPort string
		expectedRedisHost string
	}{
		{
			env:               "DEVELOPMENT",
			expectedHTTPPort:  "8080",
			expectedRedisPort: "",
			expectedRedisHost: "",
		},
		{
			env:               "KUBERNETERS-REDIS",
			expectedHTTPPort:  "8080",
			expectedRedisPort: "6379",
			expectedRedisHost: "redis-master",
		},
		{
			env:               "PRODUCTION",
			expectedHTTPPort:  "5000",
			expectedRedisPort: "",
			expectedRedisHost: "",
		},
	}
	for _, test := range tests {
		os.Setenv("CT_DNS_ENV", test.env)
		cfg := ReadConfig("../config/")
		assert.Equal(t, test.expectedHTTPPort, cfg.HTTPPort)
		assert.Equal(t, test.expectedRedisPort, cfg.Redis.Port)
		assert.Equal(t, test.expectedRedisHost, cfg.Redis.Host)
		os.Unsetenv("CT_DNS_ENV")
	}
}
