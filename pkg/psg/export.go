package psg

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/mariomac/msxmml/pkg/reader"
	"github.com/mariomac/msxmml/pkg/song"
)

const (
	tempoKey      = "tempo"
	hzKey         = "psg.hz"
	defaultBPS    = 120
	defaultHZ     = 60
	maxChannels   = 3
	defaultOctave = 4
)

type psgEncoder struct {
	bpm             int
	hz              int
	channels        channelReg
	framesCounter   int
	chFramesCounter map[string]int
	channelOrder    map[string]int
	octaves         map[string]int
}

func Export(s *song.Song) ([]byte, error) {
	// show design.md
	// reserve two bytes for the loop index
	enc, err := newPsgEncoder(s)
	if err != nil {
		return nil, err
	}
	data := make([]byte, 2, 4*1024)
	for blockNum := range s.Blocks {
		if blockNum == s.LoopIndex {
			data[0] = byte(len(data))
			data[1] = byte(len(data) >> 8)
		}
		sbr := reader.NewSyncedBlock(s.Blocks[blockNum])
		for ti, ch := sbr.Next(); ch != ""; ti, ch = sbr.Next() {
			itemData, err := enc.encodeTablatureItem(ti, ch)
			if err != nil {
				return nil, err
			}
			data = append(data, itemData...)
			if waitData := enc.encodedWaitTime(enc.nearestFrame); len(waitData) > 0 {
				data = append(data, waitData...)
			}
		}
		// At the end of a block, we need to wait for the farthest wait time
		// and sync all the channels to the current frame counter
		if waitData := enc.encodedWaitTime(enc.farthestFrame); len(waitData) > 0 {
			data = append(data, waitData...)
		}
		for k := range enc.chFramesCounter {
			enc.chFramesCounter[k] = enc.framesCounter
		}
	}
	return data, nil
}

func newPsgEncoder(s *song.Song) (*psgEncoder, error) {
	bps := defaultBPS
	hz := defaultHZ
	if bpsStr, ok := s.Properties[tempoKey]; ok {
		var err error
		if bps, err = strconv.Atoi(bpsStr); err != nil {
			return nil, fmt.Errorf("error parsing %q property: %w", tempoKey, err)
		}
	} else {
		log.Printf("assuming default %s: %v bpm", tempoKey, defaultBPS)
	}
	if hzStr, ok := s.Properties[hzKey]; ok {
		var err error
		if hz, err = strconv.Atoi(hzStr); err != nil {
			return nil, fmt.Errorf("error parsing %q property: %w", hzKey, err)
		}
	} else {
		log.Printf("assuming default %v: %v bpm", hzKey, defaultHZ)
	}
	// channel frames counter must be preloaded with all the channels
	cfc := map[string]int{}
	octaves := map[string]int{}
	for name := range s.ChannelNames {
		cfc[name] = 0
		octaves[name] = defaultOctave
	}
	return &psgEncoder{
		bpm:             bps,
		hz:              hz,
		channels:        channelReg(0b111_111), // all the channels are disabled
		framesCounter:   0,
		chFramesCounter: cfc,
		channelOrder:    map[string]int{},
		octaves:         octaves,
	}, nil
}

func (pe *psgEncoder) encodeTablatureItem(ti song.TablatureItem, channel string) ([]byte, error) {
	switch {
	case ti.Note != nil:
		instrs, err := pe.encodeNote(ti.Note, channel)
		if err != nil {
			return nil, err
		}
		return encodeInstructions(instrs), nil
	case ti.SetOctave != nil:
		pe.octaves[channel] = *ti.SetOctave
	case ti.OctaveStep != nil:
		pe.octaves[channel] += *ti.OctaveStep
	case ti.Silence != nil:
		instrs, err := pe.encodeSilence(ti.Silence, channel)
		if err != nil {
			return nil, err
		}
		return encodeInstructions(instrs), nil
	case ti.Volume != nil:
		// TODO
	case ti.Instrument != nil:
		// TODO
	default:
		panic(fmt.Sprintf("BUG! wrong value %#v", ti))
	}
	return nil, nil
}

func (pe *psgEncoder) encodedWaitTime(waitFunc func() int) []byte {
	ftw := waitFunc()
	if ftw == 0 {
		return nil
	}
	var waits []instruction

	pe.framesCounter += ftw
	// wait instruction does not allow more than 5-byte wait times (32 frames). Concatenate waits if needed
	for ftw > maxWaitValue {
		waits = append(waits, instruction{Type: wait, Data: maxWaitValue})
		ftw -= maxWaitValue
	}
	waits = append(waits, instruction{Type: wait, Data: uint16(ftw)})
	return encodeInstructions(waits)
}

func pow2(n int) int {
	pow := 1
	for i := 0; i < n; i++ {
		pow *= 2
	}
	return pow
}

func (c *psgEncoder) encodeSilence(silence *song.Silence, channel string) ([]instruction, error) {
	// enable channel, if not yet enabled
	channelOrder := c.orderFor(channel)
	if channelOrder >= maxChannels {
		return nil,
			fmt.Errorf("can't assign an order to channel %q. PSG can't handle more than 3 channels", channels)
	}
	var instrs []instruction
	// todo: make sure we enable/disable channels at the beginning of a loop, to avoid loosing status
	if c.channels.toneEnabled(channelOrder) {
		// todo: optimize: wrap multiple channel sets into one single instruction
		c.channels.disableTone(channelOrder)
		instrs = append(instrs, instruction{Type: channels, Data: uint16(c.channels)})
	}

	frames := c.framesFor(silence.Length, 1, 1)
	c.addFramesCount(channel, frames)
	return instrs, nil
}

func (c *psgEncoder) encodeNote(note *song.Note, channel string) ([]instruction, error) {
	// enable channel, if not yet enabled
	channelOrder := c.orderFor(channel)
	if channelOrder >= maxChannels {
		return nil,
			fmt.Errorf("can't assign an order to channel %q. PSG can't handle more than 3 channels", channels)
	}
	var instrs []instruction
	// todo: make sure we enable/disable channels at the beginning of a loop, to avoid loosing status
	if !c.channels.toneEnabled(channelOrder) {
		// todo: set tone/noise depending on the instrument type
		c.channels.enableTone(channelOrder)
		// todo: optimize: wrap multiple channel sets into one single instruction
		instrs = append(instrs, instruction{Type: channels, Data: uint16(c.channels)})
	}

	var noteTypes = [maxChannels]instructionType{toneA, toneB, toneC}
	// calculate how many frames we should wait after this note and advance the channel beats
	// counter
	dnd, dor := 1, 1
	for d := 1; d <= note.Dots; d++ {
		up := 1
		down := pow2(d)
		dnd = dnd*down + up*dor
		dor = dor * down
	}
	// todo: do also quatriplets, quintuplets, sextuplets, etc...
	if note.Tuplet == 3 {
		dnd *= 2
		dor *= 3
	}
	frames := c.framesFor(note.Length, dnd, dor)
	c.addFramesCount(channel, frames)

	// get tone part
	freq, err := c.frequencyFor(note, c.octaves[channel])
	if err != nil {
		return nil, err
	}
	instrs = append(instrs, instruction{
		Type: noteTypes[channelOrder],
		Data: freq,
	})
	return instrs, nil
}

// todo: calculate accumulated error
// number of frames to wait for a given length
// dividend and divisor would be 1 unless you want to alter the length size (tuplets, dots...)
func (c *psgEncoder) framesFor(length, dividend, divisor int) int {
	// 4 beats       1min      60s    c.hz frames
	// ------- * ----------- * ---- * ----------- * (dividend/divisor)
	// length    c.bpm beats   1min       1s
	return (4 * 60 * c.hz * dividend) / (length * c.bpm * divisor)
}

func (c *psgEncoder) addFramesCount(channel string, frames int) {
	if chTime, ok := c.chFramesCounter[channel]; ok {
		c.chFramesCounter[channel] = chTime + frames
	} else {
		c.chFramesCounter[channel] = frames
	}
}

func (c *psgEncoder) nearestFrame() int {
	nearest := math.MaxInt64
	for _, t := range c.chFramesCounter {
		dist := t - c.framesCounter
		if dist >= 0 && dist < nearest {
			nearest = dist
		}
	}
	return nearest
}

func (c *psgEncoder) farthestFrame() int {
	farthest := 0
	for _, t := range c.chFramesCounter {
		dist := t - c.framesCounter
		if dist > farthest {
			farthest = dist
		}
	}
	return farthest
}

func (c *psgEncoder) orderFor(channel string) int {
	if order, ok := c.channelOrder[channel]; ok {
		return order
	}
	ord := len(c.channelOrder)
	c.channelOrder[channel] = ord
	return ord
}

type noteKey struct {
	pitch  song.Pitch
	half   song.Halftone
	octave int
}

// Frequencies adapted to 12-bit PSG tone generator according to the MSX2 technical handbook
var frequencies = map[noteKey]uint16{
	{pitch: song.C, octave: 1}:                   0xD5D,
	{pitch: song.C, octave: 1, half: song.Sharp}: 0xC9C,
	{pitch: song.D, octave: 1, half: song.Flat}:  0xC9C,
	{pitch: song.D, octave: 1}:                   0xBE7,
	{pitch: song.D, octave: 1, half: song.Sharp}: 0xB3C,
	{pitch: song.E, octave: 1, half: song.Flat}:  0xB3C,
	{pitch: song.E, octave: 1}:                   0xA9B,
	{pitch: song.F, octave: 1}:                   0xA02,
	{pitch: song.F, octave: 1, half: song.Sharp}: 0x973,
	{pitch: song.G, octave: 1, half: song.Flat}:  0x973,
	{pitch: song.G, octave: 1}:                   0x8EB,
	{pitch: song.G, octave: 1, half: song.Sharp}: 0x88B,
	{pitch: song.A, octave: 1, half: song.Flat}:  0x88B,
	{pitch: song.A, octave: 1}:                   0x7F2,
	{pitch: song.A, octave: 1, half: song.Sharp}: 0x780,
	{pitch: song.B, octave: 1, half: song.Flat}:  0x780,
	{pitch: song.B, octave: 1}:                   0x714,

	{pitch: song.C, octave: 2}:                   0x6AF,
	{pitch: song.C, octave: 2, half: song.Sharp}: 0x64E,
	{pitch: song.D, octave: 2, half: song.Flat}:  0x64E,
	{pitch: song.D, octave: 2}:                   0x5F4,
	{pitch: song.D, octave: 2, half: song.Sharp}: 0x59E,
	{pitch: song.E, octave: 2, half: song.Flat}:  0x59E,
	{pitch: song.E, octave: 2}:                   0x54E,
	{pitch: song.F, octave: 2}:                   0x501,
	{pitch: song.F, octave: 2, half: song.Sharp}: 0x4BA,
	{pitch: song.G, octave: 2, half: song.Flat}:  0x4BA,
	{pitch: song.G, octave: 2}:                   0x476,
	{pitch: song.G, octave: 2, half: song.Sharp}: 0x436,
	{pitch: song.A, octave: 2, half: song.Flat}:  0x436,
	{pitch: song.A, octave: 2}:                   0x3F9,
	{pitch: song.A, octave: 2, half: song.Sharp}: 0x3C0,
	{pitch: song.B, octave: 2, half: song.Flat}:  0x3C0,
	{pitch: song.B, octave: 2}:                   0x38A,

	{pitch: song.C, octave: 3}:                   0x357,
	{pitch: song.C, octave: 3, half: song.Sharp}: 0x327,
	{pitch: song.D, octave: 3, half: song.Flat}:  0x327,
	{pitch: song.D, octave: 3}:                   0x2FA,
	{pitch: song.D, octave: 3, half: song.Sharp}: 0x2CF,
	{pitch: song.E, octave: 3, half: song.Flat}:  0x2CF,
	{pitch: song.E, octave: 3}:                   0x2A7,
	{pitch: song.F, octave: 3}:                   0x281,
	{pitch: song.F, octave: 3, half: song.Sharp}: 0x25D,
	{pitch: song.G, octave: 3, half: song.Flat}:  0x25D,
	{pitch: song.G, octave: 3}:                   0x23B,
	{pitch: song.G, octave: 3, half: song.Sharp}: 0x21B,
	{pitch: song.A, octave: 3, half: song.Flat}:  0x21B,
	{pitch: song.A, octave: 3}:                   0x1FD,
	{pitch: song.A, octave: 3, half: song.Sharp}: 0x1E0,
	{pitch: song.B, octave: 3, half: song.Flat}:  0x1E0,
	{pitch: song.B, octave: 3}:                   0x1C5,

	{pitch: song.C, octave: 4}:                   0x1AC,
	{pitch: song.C, octave: 4, half: song.Sharp}: 0x194,
	{pitch: song.D, octave: 4, half: song.Flat}:  0x194,
	{pitch: song.D, octave: 4}:                   0x17D,
	{pitch: song.D, octave: 4, half: song.Sharp}: 0x168,
	{pitch: song.E, octave: 4, half: song.Flat}:  0x168,
	{pitch: song.E, octave: 4}:                   0x153,
	{pitch: song.F, octave: 4}:                   0x140,
	{pitch: song.F, octave: 4, half: song.Sharp}: 0x12E,
	{pitch: song.G, octave: 4, half: song.Flat}:  0x12E,
	{pitch: song.G, octave: 4}:                   0x11D,
	{pitch: song.G, octave: 4, half: song.Sharp}: 0x10D,
	{pitch: song.A, octave: 4, half: song.Flat}:  0x10D,
	{pitch: song.A, octave: 4}:                   0xFE,
	{pitch: song.A, octave: 4, half: song.Sharp}: 0xF0,
	{pitch: song.B, octave: 4, half: song.Flat}:  0xF0,
	{pitch: song.B, octave: 4}:                   0xE3,

	{pitch: song.C, octave: 5}:                   0xD6,
	{pitch: song.C, octave: 5, half: song.Sharp}: 0xCA,
	{pitch: song.D, octave: 5, half: song.Flat}:  0xCA,
	{pitch: song.D, octave: 5}:                   0xBE,
	{pitch: song.D, octave: 5, half: song.Sharp}: 0x84,
	{pitch: song.E, octave: 5, half: song.Flat}:  0x84,
	{pitch: song.E, octave: 5}:                   0xAA,
	{pitch: song.F, octave: 5}:                   0xA0,
	{pitch: song.F, octave: 5, half: song.Sharp}: 0x97,
	{pitch: song.G, octave: 5, half: song.Flat}:  0x97,
	{pitch: song.G, octave: 5}:                   0x8F,
	{pitch: song.G, octave: 5, half: song.Sharp}: 0x87,
	{pitch: song.A, octave: 5, half: song.Flat}:  0x87,
	{pitch: song.A, octave: 5}:                   0x7F,
	{pitch: song.A, octave: 5, half: song.Sharp}: 0x78,
	{pitch: song.B, octave: 5, half: song.Flat}:  0x78,
	{pitch: song.B, octave: 5}:                   0x71,

	{pitch: song.C, octave: 6}:                   0x6B,
	{pitch: song.C, octave: 6, half: song.Sharp}: 0x65,
	{pitch: song.D, octave: 6, half: song.Flat}:  0x65,
	{pitch: song.D, octave: 6}:                   0x5F,
	{pitch: song.D, octave: 6, half: song.Sharp}: 0x5A,
	{pitch: song.E, octave: 6, half: song.Flat}:  0x5A,
	{pitch: song.E, octave: 6}:                   0x55,
	{pitch: song.F, octave: 6}:                   0x50,
	{pitch: song.F, octave: 6, half: song.Sharp}: 0x4C,
	{pitch: song.G, octave: 6, half: song.Flat}:  0x4C,
	{pitch: song.G, octave: 6}:                   0x47,
	{pitch: song.G, octave: 6, half: song.Sharp}: 0x43,
	{pitch: song.A, octave: 6, half: song.Flat}:  0x43,
	{pitch: song.A, octave: 6}:                   0x40,
	{pitch: song.A, octave: 6, half: song.Sharp}: 0x3C,
	{pitch: song.B, octave: 6, half: song.Flat}:  0x3C,
	{pitch: song.B, octave: 6}:                   0x39,

	{pitch: song.C, octave: 7}:                   0x35,
	{pitch: song.C, octave: 7, half: song.Sharp}: 0x32,
	{pitch: song.D, octave: 7, half: song.Flat}:  0x32,
	{pitch: song.D, octave: 7}:                   0x30,
	{pitch: song.D, octave: 7, half: song.Sharp}: 0x2D,
	{pitch: song.E, octave: 7, half: song.Flat}:  0x2D,
	{pitch: song.E, octave: 7}:                   0x2A,
	{pitch: song.F, octave: 7}:                   0x28,
	{pitch: song.F, octave: 7, half: song.Sharp}: 0x26,
	{pitch: song.G, octave: 7, half: song.Flat}:  0x26,
	{pitch: song.G, octave: 7}:                   0x24,
	{pitch: song.G, octave: 7, half: song.Sharp}: 0x22,
	{pitch: song.A, octave: 7, half: song.Flat}:  0x22,
	{pitch: song.A, octave: 7}:                   0x20,
	{pitch: song.A, octave: 7, half: song.Sharp}: 0x1E,
	{pitch: song.B, octave: 7, half: song.Flat}:  0x1E,
	{pitch: song.B, octave: 7}:                   0x1C,

	{pitch: song.C, octave: 8}:                   0x1B,
	{pitch: song.C, octave: 8, half: song.Sharp}: 0x19,
	{pitch: song.D, octave: 8, half: song.Flat}:  0x19,
	{pitch: song.D, octave: 8}:                   0x18,
	{pitch: song.D, octave: 8, half: song.Sharp}: 0x16,
	{pitch: song.E, octave: 8, half: song.Flat}:  0x16,
	{pitch: song.E, octave: 8}:                   0x15,
	{pitch: song.F, octave: 8}:                   0x14,
	{pitch: song.F, octave: 8, half: song.Sharp}: 0x13,
	{pitch: song.G, octave: 8, half: song.Flat}:  0x13,
	{pitch: song.G, octave: 8}:                   0x12,
	{pitch: song.G, octave: 8, half: song.Sharp}: 0x11,
	{pitch: song.A, octave: 8, half: song.Flat}:  0x11,
	{pitch: song.A, octave: 8}:                   0x10,
	{pitch: song.A, octave: 8, half: song.Sharp}: 0xF,
	{pitch: song.B, octave: 8, half: song.Flat}:  0xF,
	{pitch: song.B, octave: 8}:                   0xE,
}

func (c *psgEncoder) frequencyFor(n *song.Note, octave int) (uint16, error) {
	freq, ok := frequencies[noteKey{pitch: n.Pitch, half: n.Halftone, octave: octave}]
	if !ok {
		ht := byte(n.Halftone)
		if ht == 0 {
			ht = ' '
		}
		return 0, fmt.Errorf("unsupported note: %c%c for octave %d", n.Pitch, ht, octave)
	}
	return freq, nil
}
