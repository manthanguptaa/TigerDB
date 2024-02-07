package main

import (
	"TigerDB/cache"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "listen address of the server")
	leaderAddr := flag.String("leaderaddr", "", "listen address of the leader")
	flag.Parse()

	opts := ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   true,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte("SET Foo Bar 2600000000"))
		time.Sleep(time.Second * 2)

		conn.Write([]byte("GET Foo"))
		buf := make([]byte, 1000)
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))
	}()

	server := NewServer(opts, cache.NewCache())
	server.Start()
}
