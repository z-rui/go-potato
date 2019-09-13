package potato

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

type Player struct {
	RingMaster Address

	rand   *rand.Rand
	n      int
	id     PlayerID
	rpc_rm *rpc.Client
	rpc_p  [2]*rpc.Client
}

func (p *Player) Run(network, address string) {
	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	server := rpc.NewServer()
	server.Register(p)
	p.rpc_rm, err = rpc.Dial(p.RingMaster.Network, p.RingMaster.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer p.rpc_rm.Close()
	var rinfo RingInfo
	err = p.rpc_rm.Call("RingMaster.EnterRing", Address{network, address}, &rinfo)
	if err != nil {
		log.Fatal(err)
	}
	p.n = rinfo.N
	p.id = rinfo.ID
	log.Printf("Connected as player %d", p.id)
	p.rand = rand.New(rand.NewSource(time.Now().Unix() + int64(p.id)))
	p.rpc_p[0], err = rpc.Dial(rinfo.L.Network, rinfo.L.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer p.rpc_p[0].Close()
	p.rpc_p[1], err = rpc.Dial(rinfo.R.Network, rinfo.R.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer p.rpc_p[1].Close()
	p.rpc_p[1], err = rpc.Dial(rinfo.R.Network, rinfo.R.Address)
	go server.Accept(ln)
	p.rpc_rm.Call("RingMaster.WaitUntilEnd", struct{}{}, nil)
}

func (p *Player) Receive(potato *Potato, _ *struct{}) error {
	if potato.TTL <= 0 {
		return fmt.Errorf("malformed potato: TTL = %d", potato.TTL)
	}
	potato.TTL--
	potato.Trace = append(potato.Trace, p.id)
	go p.sendPotato(potato)
	return nil
}

func (p *Player) sendPotato(potato *Potato) {
	var client *rpc.Client
	var call string
	if potato.TTL > 0 {
		var nextID PlayerID
		k := p.rand.Intn(2)
		client = p.rpc_p[k]
		nextID = PlayerID((int(p.id) + p.n + k + k - 1) % p.n)
		log.Printf("Sending potato to player %d.", nextID)
		call = "Player.Receive"
	} else {
		log.Print("I'm it.")
		client = p.rpc_rm
		call = "RingMaster.Receive"
	}
	err := client.Call(call, potato, nil)
	if err != nil {
		log.Fatal(err)
	}
}
