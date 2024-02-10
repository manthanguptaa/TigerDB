package main

import (
	"TigerDB/cache"
	"TigerDB/client"
	"context"
	"flag"
	"fmt"
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
		time.Sleep(time.Second * 1)
		SendStuff()
	}()

	server := NewServer(opts, cache.NewCache())
	server.Start()
}

func SendStuff() {
	for i := 0; i < 100; i++ {
		go func(i int) {
			client, err := client.New(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}
			var key = fmt.Sprintf("key_%d", i)
			var value = fmt.Sprintf("value_%d", i)

			err = client.Set(context.Background(), []byte(key), []byte(value), 0)
			if err != nil {
				log.Println(err)
			}

			resp, err := client.Get(context.Background(), []byte(key))

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(resp))
			client.Close()
		}(i)
	}

}
