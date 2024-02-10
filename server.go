package main

import (
	"TigerDB/cache"
	"TigerDB/proto"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
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
	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
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
