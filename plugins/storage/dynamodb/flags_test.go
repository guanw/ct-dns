package dynamodb

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFlags(t *testing.T) {
	flagSet := flag.NewFlagSet("dynamodb", flag.ExitOnError)
	AddFlags(flagSet)
	flagSet.Parse([]string{"--dynamodb-region", "us-east-2"})
	if flagSet.Parsed() {
		assert.Equal(t, flagSet.Lookup("dynamodb-region").Value.String(), "us-east-2")
		assert.Equal(t, flagSet.Lookup("dynamodb-endpoint").Value.String(), "http://localhost:8000")
	}
}
