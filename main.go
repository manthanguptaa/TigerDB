package main

import (
	"TigerDB/cache"
	"TigerDB/client"
	"context"
	"flag"
	"log"
	"time"
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

	go func() {
		time.Sleep(time.Second * 2)
		client, err := client.New(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 5; i++ {
			SendCommand(client)
			time.Sleep(time.Millisecond * 200)
		}
		client.Close()
		time.Sleep(time.Millisecond * 200)
	}()

	server := NewServer(opts, cache.NewCache())
	server.Start()
}

func SendCommand(c *client.Client) {
	_, err := c.Set(context.Background(), []byte("Manthan"), []byte("Gupta"), 0)
	if err != nil {
		log.Fatal(err)
	}

}
