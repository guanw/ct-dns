package dynamodb

import (
	"flag"
)

func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String("dynamodb-region", "us-east-1", "--dynamodb-region <name>")
	flagSet.String("dynamodb-endpoint", "http://localhost:8000", "--dynamodb-endpoint <endpoint>")
}
