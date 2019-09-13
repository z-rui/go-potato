#!/bin/sh

if [ $# -ne 2 ]; then
	echo "Syntax: $0 <num_players> <ttl>"
	exit 1
fi

NUM_PLAYERS=$1
TTL=$2
PORT_BASE=10000
RM_PORT=23333

ringmaster/ringmaster -r ":$RM_PORT" -ttl $TTL -n $NUM_PLAYERS &
sleep 1  # wait for the server to set up
i=0
while [ $i -lt $NUM_PLAYERS ]; do
	player/player -r ":$RM_PORT" -l ":$((PORT_BASE+i))" &
	i=$((i+1))
done

wait
