package api

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemCpy(t *testing.T) {
	dest := []byte{1, 2, 3, 4}
	src := []byte{5, 6, 7, 8}
	memcpy(dest, src, 2)
	fmt.Printf("%v", dest)
}

func TestMemSet(t *testing.T) {
	length := 6
	dest := bytes.Repeat([]byte{0}, length)
	memset(dest, byte(1), 4)
	assert.Equal(t, []byte{1, 1, 1, 1, 0, 0}, dest)
	dest = bytes.Repeat([]byte{0}, length)
	memset(dest, byte(1), 6)
	assert.Equal(t, []byte{1, 1, 1, 1, 1, 1}, dest)
	fmt.Printf("%v", string(dest))
}

func TestMemCmp(t *testing.T) {
	buffer1 := []byte("DWgaOtP12df0")
	buffer2 := []byte("DWGAOTP12DF0")
	n := memcmp(buffer1, buffer2, uint32(len(buffer1)))
	assert.Equal(t, int32(1), n)
}
