# Hot potato!

This is an implementation of the hot potato game using Go's `net/rpc` package.

## Overview

Each player is a rpc server providing Receive to receive the potato.
The player is also a rpc client to the ring master and its two neighbors.

The ring master is a rpc server providing the following methods:
1. EnterRing.  Player calls this method with its server address.
   The ring master returns information about the ring back to the player.
2. WaitUntilEnd.  Player calls this method to be blocked until the end of
   the game.
3. Receive.  The ring master receives the potato from the last player.

With `net/rpc`, not effort is spent on (de)serializing the data structures
or network operations.
