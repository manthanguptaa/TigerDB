package main

import (
	"TigerDB/client"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/raft"
)

func main() {
	var (
		cfg      = raft.DefaultConfig()
		fsm      = &raft.MockFSM{}
		logStore = raft.NewInmemStore()
		timeout  = time.Second * 5
		stable   = raft.NewInmemStore()
	)

	snapshotStore, err := raft.NewFileSnapshotStore("log", 3, nil)
	if err != nil {
		log.Fatal(err)
	}

	cfg.LocalID = "RandomString"

	tr, err := raft.NewTCPTransport("localhost:4000", nil, 10, timeout, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	follower_1 := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID(cfg.LocalID),
		Address:  raft.ServerAddress("localhost:4000"),
	}
	follower_2 := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID("follower_2"),
		Address:  raft.ServerAddress("localhost:4001"),
	}
	follower_3 := raft.Server{
		Suffrage: raft.Voter,
		ID:       raft.ServerID("follower_3"),
		Address:  raft.ServerAddress("localhost:4002"),
	}

	serverConfig := raft.Configuration{
		Servers: []raft.Server{follower_1, follower_2, follower_3},
	}

	r, err := raft.NewRaft(cfg, fsm, logStore, stable, snapshotStore, tr)
	if err != nil {
		log.Fatal(err)
	}

	r.BootstrapCluster(serverConfig)

	fmt.Printf("%+v\n", r)

	select {}
}

func SendStuff() {
	c, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		var key = fmt.Sprintf("key_%d", i)
		var value = fmt.Sprintf("value_%d", i)

		err = c.Set(context.Background(), []byte(key), []byte(value), 0)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second)
	}
	c.Close()
}
