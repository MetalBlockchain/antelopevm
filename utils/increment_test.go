package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrement(t *testing.T) {
	increment := Increment(0)
	increment.Increment()
	assert.Equal(t, increment, Increment(1), "Increment should be 1")
}
