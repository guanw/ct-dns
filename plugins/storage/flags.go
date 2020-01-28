package storage

import "flag"

// AddFlags add flags to for storage initialization
func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String("storage-type", "redis", "--storage-type <name>")
}
