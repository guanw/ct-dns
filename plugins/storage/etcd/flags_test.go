package etcd

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFlags(t *testing.T) {
	flagSet := flag.NewFlagSet("etcd", flag.ExitOnError)
	AddFlags(flagSet)
	flagSet.Parse([]string{"--etcd-endpoints", "http://192.0.0.1:5000,http://192.0.0.1:5001"})
	if flagSet.Parsed() {
		assert.Equal(t, flagSet.Lookup("etcd-endpoints").Value.String(), "http://192.0.0.1:5000,http://192.0.0.1:5001")
	}
}
