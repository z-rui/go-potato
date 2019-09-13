package main

import (
	"flag"

	"github.com/z-rui/go-potato"
)

var (
	ringmaster = flag.String("r", ":2333", "address of ring master")
	player     = flag.String("l", ":10000", "address of player")
)

func main() {
	flag.Parse()
	p := &potato.Player{
		RingMaster: potato.Address{
			Network: "tcp",
			Address: *ringmaster,
		},
	}
	p.Run("tcp", *player)
}
