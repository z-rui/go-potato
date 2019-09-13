package potato

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type RingMaster struct {
	N   int
	TTL int

	player  []Address
	mtx     sync.Mutex
	wg_conn sync.WaitGroup
	wg_game sync.WaitGroup
}

func (rm *RingMaster) Run(network, address string) {
	rm.player = nil
	rm.wg_conn.Add(rm.N)

	ln, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	server := rpc.NewServer()
	server.Register(rm)
	go server.Accept(ln)
	rm.wg_conn.Wait()
	if len(rm.player) > 0 && rm.TTL > 0 {
		rm.wg_game.Add(1)
		rm.sendPotato()
	}
	rm.wg_game.Wait()
}

func (rm *RingMaster) EnterRing(playerAddr Address, resp *RingInfo) error {
	id, err := rm.addPlayer(playerAddr)
	if err != nil {
		return err
	}
	rm.wg_conn.Wait() // wait for N connections
	resp.N = rm.N
	resp.ID = id
	resp.L = rm.player[(int(id)+rm.N-1)%rm.N]
	resp.R = rm.player[(int(id)+1)%rm.N]
	return nil
}

func (rm *RingMaster) WaitUntilEnd(_ struct{}, _ *struct{}) error {
	rm.wg_game.Wait()
	return nil
}

func (rm *RingMaster) Receive(potato Potato, _ *struct{}) error {
	log.Printf("Received potato, trace = %v", potato.Trace)
	rm.wg_game.Done()
	return nil
}

func (rm *RingMaster) addPlayer(playerAddr Address) (id PlayerID, err error) {
	rm.mtx.Lock()
	defer rm.mtx.Unlock()
	id = PlayerID(len(rm.player))
	if int(id) >= rm.N {
		err = fmt.Errorf("There are already %d players in the game", id)
		return
	}
	rm.player = append(rm.player, playerAddr)
	log.Printf("Player %d connected", id)
	rm.wg_conn.Done()
	return
}

func (rm *RingMaster) sendPotato() {
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	id := PlayerID(rng.Intn(rm.N))
	log.Printf("Sending potato to player %d", id)
	p := rm.player[id]
	client, err := rpc.Dial(p.Network, p.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	potato := &Potato{TTL: rm.TTL}
	err = client.Call("Player.Receive", potato, nil)
	if err != nil {
		log.Fatal(err)
	}
}
