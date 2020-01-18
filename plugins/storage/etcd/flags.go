package etcd

import (
	"flag"
)

// AddFlags binds flags to etcd setup
func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String("etcd-endpoints", `http://10.110.251.205:2379`, "--etcd-endpoints <name>")
}
