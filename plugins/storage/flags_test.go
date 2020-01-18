package storage

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddFlags(t *testing.T) {
	flagSet := flag.NewFlagSet("storage", flag.ExitOnError)
	AddFlags(flagSet)
	flagSet.Parse([]string{"--storage-type", "memory"})
	if flagSet.Parsed() {
		assert.Equal(t, flagSet.Lookup("storage-type").Value.String(), "memory")
	}
}
