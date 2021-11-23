package psg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelsRegister_Enable(t *testing.T) {
	ch := channelReg(0b111_111)
	for i := 0; i < 3; i++ {
		assert.False(t, ch.toneEnabled(i))
		assert.False(t, ch.noiseEnabled(i))
	}

	ch.enableTone(0)
	assert.True(t, ch.toneEnabled(0))
	assert.Equal(t, channelReg(0b111_110), ch)
	ch.enableNoise(0)
	assert.True(t, ch.noiseEnabled(0))
	assert.Equal(t, channelReg(0b110_110), ch)
	ch.enableTone(1)
	assert.True(t, ch.toneEnabled(1))
	assert.Equal(t, channelReg(0b110_100), ch)
	ch.enableNoise(1)
	assert.True(t, ch.noiseEnabled(1))
	assert.Equal(t, channelReg(0b100_100), ch)
	ch.enableTone(2)
	assert.True(t, ch.toneEnabled(2))
	assert.Equal(t, channelReg(0b100_000), ch)
	ch.enableNoise(2)
	assert.True(t, ch.noiseEnabled(2))
	assert.Equal(t, channelReg(0b000_000), ch)
}

func TestChannelsRegister_Disable(t *testing.T) {
	ch := channelReg(0)
	for i := 0; i < 3; i++ {
		assert.True(t, ch.toneEnabled(i))
		assert.True(t, ch.noiseEnabled(i))
	}
	ch.disableTone(0)
	assert.False(t, ch.toneEnabled(0))
	assert.Equal(t, channelReg(0b000_001), ch)
	ch.disableNoise(0)
	assert.False(t, ch.noiseEnabled(0))
	assert.Equal(t, channelReg(0b001_001), ch)
	ch.disableTone(1)
	assert.False(t, ch.toneEnabled(1))
	assert.Equal(t, channelReg(0b001_011), ch)
	ch.disableNoise(1)
	assert.False(t, ch.noiseEnabled(1))
	assert.Equal(t, channelReg(0b011_011), ch)
	ch.disableTone(2)
	assert.False(t, ch.toneEnabled(2))
	assert.Equal(t, channelReg(0b011_111), ch)
	ch.disableNoise(2)
	assert.False(t, ch.noiseEnabled(2))
	assert.Equal(t, channelReg(0b111_111), ch)
}
