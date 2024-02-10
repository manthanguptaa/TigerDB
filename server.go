package main

import (
	"TigerDB/cache"
	"TigerDB/client"
	"TigerDB/proto"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"go.uber.org/zap"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	members map[*client.Client]struct{}
	cache   cache.Cacher
	logger  *zap.SugaredLogger
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	l, _ := zap.NewProduction()
	lsugar := l.Sugar()
	return &Server{
		ServerOpts: opts,
		cache:      c,
		members:    make(map[*client.Client]struct{}),
		logger:     lsugar,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
	}

	s.logger.Infow("server starting", "addr", s.ListenAddr, "leader", s.IsLeader)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("faild to dial the leader: [%s]", s.LeaderAddr)
	}

	s.logger.Infow("connected to leader", "addr", s.LeaderAddr)

	binary.Write(conn, binary.LittleEndian, proto.CmdJoin)

	s.handleConn(conn)
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error: ", err)
			break
		}
		go s.handleCommand(conn, cmd)
	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
		s.handleGetCommand(conn, v)
	case *proto.CommandJoin:
		s.handleJoinCommand(conn, v)
	}
}

func (s *Server) handleJoinCommand(conn net.Conn, cmd *proto.CommandJoin) error {
	fmt.Println("member just joined the cluster: ", conn.RemoteAddr())
	s.members[client.NewClientFromConn(conn)] = struct{}{}
	return nil
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	log.Printf("SET %s to %s", cmd.Key, cmd.Value)

	go func() {
		for member := range s.members {
			if err := member.Set(context.TODO(), cmd.Key, cmd.Value, cmd.TTL); err != nil {
				log.Println("write failed on a follower: ", err)
			}
		}
	}()

	resp := proto.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = proto.StatusError
		serialize, _ := resp.Bytes()
		_, err := conn.Write(serialize)
		return err
	}

	resp.Status = proto.StatusOK
	serialize, _ := resp.Bytes()
	_, err := conn.Write(serialize)
	return err
}

func (s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {

	resp := proto.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusKeyNotFound
		serialize, err := resp.Bytes()
		if err != nil {
			return err
		}
		_, err = conn.Write(serialize)
		return err
	}

	resp.Status = proto.StatusOK
	resp.Value = value
	serialize, err := resp.Bytes()
	if err != nil {
		return err
	}
	_, err = conn.Write(serialize)
	if err != nil {
		return err
	}
	return nil
}
