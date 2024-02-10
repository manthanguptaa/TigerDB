package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command byte

type Status byte

func (s Status) String() string {
	switch s {
	case StatusError:
		return "Err"
	case StatusOK:
		return "OK"
	case StatusKeyNotFound:
		return "KEY NOT FOUND"
	default:
		return "NONE"
	}
}

const (
	StatusNone Status = iota
	StatusOK
	StatusError
	StatusKeyNotFound
)

const (
	CmdNonce Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

type ResponseSet struct {
	Status Status
}

type ResponseGet struct {
	Status Status
	Value  []byte
}

type CommandJoin struct {
}

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

type CommandGet struct {
	Key []byte
}

func (r *ResponseGet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}

	valueLen := int32(len(r.Value))
	if err := binary.Write(buf, binary.LittleEndian, valueLen); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, r.Value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *ResponseSet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, r.Status); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *CommandGet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, CmdGet); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, int32(len(c.Key))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Key); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *CommandSet) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, CmdSet); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, int32(len(c.Key))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Key); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, int32(len(c.Value))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, c.Value); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, int32(c.TTL)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ParseSetResponse(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}
	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return nil, err
	}
	return resp, nil
}

func ParseGetResponse(r io.Reader) (*ResponseGet, error) {
	resp := &ResponseGet{}
	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return nil, err
	}

	var valueLen int32
	if err := binary.Read(r, binary.LittleEndian, &valueLen); err != nil {
		return nil, err
	}

	resp.Value = make([]byte, valueLen)
	if err := binary.Read(r, binary.LittleEndian, &resp.Value); err != nil {
		return nil, err
	}

	return resp, nil
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case CmdSet:
		set, err := parseSetCommand(r)
		return set, err
	case CmdGet:
		get, err := parseGetCommand(r)
		return get, err
	case CmdJoin:
		return &CommandJoin{}, nil
	default:
		return nil, fmt.Errorf("invalid command")
	}
}

func parseSetCommand(r io.Reader) (*CommandSet, error) {
	cmd := &CommandSet{}

	var keyLen int32
	if err := binary.Read(r, binary.LittleEndian, &keyLen); err != nil {
		return nil, err
	}
	cmd.Key = make([]byte, keyLen)
	if err := binary.Read(r, binary.LittleEndian, &cmd.Key); err != nil {
		return nil, err
	}

	var valueLen int32
	if err := binary.Read(r, binary.LittleEndian, &valueLen); err != nil {
		return nil, err
	}
	cmd.Value = make([]byte, valueLen)
	if err := binary.Read(r, binary.LittleEndian, &cmd.Value); err != nil {
		return nil, err
	}

	var ttl int32
	if err := binary.Read(r, binary.LittleEndian, &ttl); err != nil {
		return nil, err
	}
	cmd.TTL = int(ttl)

	return cmd, nil
}

func parseGetCommand(r io.Reader) (*CommandGet, error) {
	cmd := &CommandGet{}

	var keyLen int32
	if err := binary.Read(r, binary.LittleEndian, &keyLen); err != nil {
		return nil, err
	}
	cmd.Key = make([]byte, keyLen)
	if err := binary.Read(r, binary.LittleEndian, &cmd.Key); err != nil {
		return nil, err
	}

	return cmd, nil
}
