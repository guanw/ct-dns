package redis

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFlags(t *testing.T) {
	flagSet := flag.NewFlagSet("redis", flag.ExitOnError)
	AddFlags(flagSet)
	flagSet.Parse([]string{"--redis-endpoint", "0.0.0.0:6379"})
	if flagSet.Parsed() {
		assert.Equal(t, flagSet.Lookup("redis-endpoint").Value.String(), "0.0.0.0:6379")
	}
}
