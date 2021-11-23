package psg

type channelReg uint8

func (c *channelReg) enableTone(channelOrder int) {
	*c &= ^(1 << channelOrder)
}

func (c *channelReg) disableTone(channelOrder int) {
	*c |= 1 << channelOrder
}

func (c *channelReg) enableNoise(channelOrder int) {
	*c &= ^(0b1000 << channelOrder)
}

func (c *channelReg) disableNoise(channelOrder int) {
	*c |= 0b1000 << channelOrder
}

func (c *channelReg) toneEnabled(channelOrder int) bool {
	return *c&(1<<channelOrder) == 0
}

func (c *channelReg) noiseEnabled(channelOrder int) bool {
	return *c&(0b1000<<channelOrder) == 0
}
