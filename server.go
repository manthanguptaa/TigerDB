package main

import (
	"TigerDB/cache"
	"context"
	"fmt"
	"log"
	"net"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	followers map[net.Conn]struct{}
	cache     cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		followers:  make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("conn read error: %s\n", err)
			break
		}

		go s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("failed to parse command: ", err)
		conn.Write([]byte(err.Error()))
		return
	}

	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCmd(conn, msg)
	case CMDGet:
		_, err = s.handleGetCmd(conn, msg)
	}

	if err != nil {
		fmt.Println("failed to handle command: ", err)
		conn.Write([]byte(err.Error()))
	}
}

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) ([]byte, error) {
	value, err := s.cache.Get(msg.Key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)
	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}
