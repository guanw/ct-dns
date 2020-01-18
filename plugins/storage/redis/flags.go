package redis

import "flag"

// AddFlags binds flags to redis setup
func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String("redis-endpoint", "0.0.0.0:6379", "--redis-endpoint <name>")
}
