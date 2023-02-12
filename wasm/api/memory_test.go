package api

import (
	"fmt"
	"testing"
)

func TestMemCpy(t *testing.T) {
	dest := []byte{1, 2, 3, 4}
	src := []byte{5, 6, 7, 8}
	memcpy(dest, src, 2)
	fmt.Printf("%v", dest)
}

func TestMemSet(t *testing.T) {
	dest := []byte("This is string.h library function")
	memset(dest, []byte(".")[0], 4)
	fmt.Printf("%v", string(dest))
}
