package proto

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetCommand(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	serialized, err := cmd.Bytes()
	assert.Nil(t, err)

	fmt.Println(serialized)

	r := bytes.NewReader(serialized)

	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)
	assert.Equal(t, cmd, pcmd)
}

func TestParseGetCommand(t *testing.T) {
	cmd := &CommandGet{
		Key: []byte("Foo"),
	}

	serialized, err := cmd.Bytes()
	assert.Nil(t, err)

	r := bytes.NewReader(serialized)

	pcmd, err := ParseCommand(r)
	assert.Nil(t, err)
	assert.Equal(t, cmd, pcmd)
}

func BenchmarkParseCommand(b *testing.B) {
	cmd := &CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}
	for i := 0; i < b.N; i++ {
		serialized, _ := cmd.Bytes()
		r := bytes.NewReader(serialized)
		ParseCommand(r)
	}
}
