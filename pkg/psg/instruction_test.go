package psg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeInstructions(t *testing.T) {
	encoded := encodeInstructions([]instruction{
		{Type: envelopeCycle, Data: 0xABCD},
		{Type: wait, Data: 0b10101},
		{Type: toneA, Data: 0xDCA},
		{Type: toneB, Data: 0xADC},
		{Type: toneC, Data: 0x123},
		{Type: envelopeShape, Data: 0b1010},
		{Type: noiseRate, Data: 0b11000},
		{Type: channels, Data: 0b010101},
		{Type: volumeA, Data: 0b1111},
		{Type: volumeB, Data: 0b1001},
		{Type: volumeC, Data: 0b1110},
		{Type: envelopeA},
		{Type: envelopeB},
		{Type: envelopeC},
		{Type: end},
	})
	assert.Equal(t, []byte{
		0, 0xAB, 0xCD,
		0b10101,
		0x2D, 0xCA,
		0x3A, 0xDC,
		0x71, 0x23,
		0b01101010,
		0b01011000,
		0b10010101,
		0b11001111,
		0b11011001,
		0b11101110,
		0b11110000,
		0b11110001,
		0b11110010,
		0b11111000,
	}, encoded)
}
