// Potato is a package that demonstrates the functionalities of Go's
// net/rpc package by implementing the logic of the "hot-potato" game.
//
// The "hot-potato" game is played by a ring master and N players.
// They communicate over the network, passing a "potato" object.
// The potato has a TTL (time-to-alive), which specifies how many turns it
// will be passed from one to another, and a trace, which records which players
// the potato has been passed to.
//
// The overall process of the game is:
//   1.  The ring master starts, knowing N and TTL.
//   2.  Each player connects to the ring master.  The ring master will tell
//       the player its ID and the network address of its neighbors.
//       Each player has 2 neighbors, and the players form a ring.
//   3.  Once all players has connected to the ring master, the game starts.
//       The ring master randomly chooses a player and sends the potato to it.
//   4.  Upon receiving the potato, the player decrements its TTL and appends
//       the player ID to the trace.
//   5a. If TTL>0, then the player randomly chooses a neighbor and sends
//       the potato to it, and prints a message describing where is potato was
//       sent.
//   5b. If TTL=0, then the player prints a message, "I'm it.", and sends
//       the potato to the ring master.
//   6.  Once the ring master receives the potato, the game ends.  The ring
//       master should print the trace and notify all players to shut down.
//
// The player can assume that if its ID is i, then the "left" neighbor has
// ID (i+N-1)%N and the "right" neighbor has ID (i+1)%N.
//
// For the special case where TTL=0 and the start of the game, the ring master
// should immediately shut down the game and not send the potato to any player.
package potato

// PlayerID is an integer for identifying the player.
// For a hot-potato game of N players, the IDs are numbered from 0 to N-1.
type PlayerID int

// Potato is the object to be passed between players,
// or between a player and the ring master at the start and end of the game.
type Potato struct {
	TTL   int
	Trace []PlayerID
}

type Address struct {
	Network string
	Address string
}

type RingInfo struct {
	N    int      // total number of players
	ID   PlayerID // ID assigned by ring master
	L, R Address  // address of neighbors
}
