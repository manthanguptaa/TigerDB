package main

import (
	"TigerDB/cache"
	"flag"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen address of the server")
	leaderAddr := flag.String("leaderaddr", "", "listen address of the leader")
	flag.Parse()

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	server := NewServer(opts, cache.NewCache())
	server.Start()
}
