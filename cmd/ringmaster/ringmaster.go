package main

import (
	"flag"

	"github.com/z-rui/go-potato"
)

var (
	n          = flag.Int("n", 0, "number of players")
	ttl        = flag.Int("ttl", 0, "number of hops")
	ringmaster = flag.String("r", ":2333", "address of ring master")
)

func main() {
	flag.Parse()
	rm := &potato.RingMaster{
		N:   *n,
		TTL: *ttl,
	}
	rm.Run("tcp", *ringmaster)
}
