package psg

import "fmt"

type instructionType int

const (
	envelopeCycle instructionType = iota
	wait
	toneA
	toneB
	toneC
	envelopeShape
	noiseRate
	channels
	volumeA
	volumeB
	volumeC
	envelopeA
	envelopeB
	envelopeC
)

type instruction struct {
	Type instructionType
	Data uint16
}

func (i *instruction) encode() []byte {
	switch i.Type {
	case envelopeCycle:
		return []byte{0, byte(i.Data >> 8), byte(i.Data)}
	case wait:
		return []byte{byte(i.Data & 0b11111)}
	case toneA:
		return []byte{0b00100000 | byte((i.Data>>8)&0b1111), byte(i.Data)}
	case toneB:
		return []byte{0b00110000 | byte((i.Data>>8)&0b1111), byte(i.Data)}
	case toneC:
		return []byte{0b01110000 | byte((i.Data>>8)&0b1111), byte(i.Data)}
	case envelopeShape:
		return []byte{0b01100000 | byte(i.Data&0b1111)}
	case noiseRate:
		return []byte{0b01000000 | byte(i.Data&0b11111)}
	case channels:
		return []byte{0b10000000 | byte(i.Data&0b111111)}
	case volumeA:
		return []byte{0b11000000 | byte(i.Data&0b1111)}
	case volumeB:
		return []byte{0b11010000 | byte(i.Data&0b1111)}
	case volumeC:
		return []byte{0b11100000 | byte(i.Data&0b1111)}
	case envelopeA:
		return []byte{0b11110000}
	case envelopeB:
		return []byte{0b11110001}
	case envelopeC:
		return []byte{0b11110010}
	}
	panic(fmt.Sprintf("Unknown instruction type: %d", i))
}

func encodeInstructions(ints []instruction) []byte {
	var encoded []byte
	for _, i := range ints {
		encoded = append(encoded, i.encode()...)
	}
	return encoded
}