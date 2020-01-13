package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignExtend(t *testing.T) {
	assert.Equal(t, uint16(0b1111_1111_1111_1111), SignExtend(0b11111, 5))
	assert.Equal(t, uint16(0b0000_0000_0000_1111), SignExtend(0b01111, 5))
}
