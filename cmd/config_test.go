package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadConfig(t *testing.T) {
	tests := []struct {
		env          string
		expectedPort string
	}{
		{
			env:          "DEVELOPMENT",
			expectedPort: "8080",
		},
		{
			env:          "PRODUCTION",
			expectedPort: "5000",
		},
	}
	for _, test := range tests {
		os.Setenv("CT_DNS_ENV", test.env)
		cfg := ReadConfig("../config/")
		assert.Equal(t, test.expectedPort, cfg.HTTPPort)
		os.Unsetenv("CT_DNS_ENV")
	}
}
