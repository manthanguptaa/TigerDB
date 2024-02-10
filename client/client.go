package client

import (
	"TigerDB/proto"
	"context"
	"fmt"
	"net"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func NewClientFromConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet{
		Key: key,
	}

	serialize, _ := cmd.Bytes()
	_, err := c.conn.Write(serialize)
	if err != nil {
		return nil, err
	}

	resp, err := proto.ParseGetResponse(c.conn)
	if err != nil {
		return nil, err
	}

	if resp.Status == proto.StatusKeyNotFound {
		return nil, fmt.Errorf("couldn't find the key [%s]", key)
	}

	if resp.Status != proto.StatusOK {
		return nil, fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}

	return resp.Value, nil
}

func (c *Client) Set(ctx context.Context, key []byte, value []byte, ttl int) error {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	serialize, _ := cmd.Bytes()
	_, err := c.conn.Write(serialize)
	if err != nil {
		return err
	}
	resp, err := proto.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != proto.StatusOK {
		return fmt.Errorf("server responded with non OK status [%s]", resp.Status)
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
