package bots

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExample1(t *testing.T) {
	assert.InDelta(t, 0, getExample1().Input[10], 0.05)
	assert.InDelta(t, 1, getExample1().Input[11], 0.05)
	assert.InDelta(t, 1, getExample1().Input[12], 0.05)
	assert.InDelta(t, 0, getExample1().Response[12], 0.05)
	assert.InDelta(t, 1, getExample1().Response[13], 0.05)
}
