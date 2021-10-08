package lang

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mariomac/msxmml/pkg/song"
)

const (
	defaultOctave = 4
	minOctave     = 0
	maxOctave     = 8
	minLength     = 1
	maxLength     = 64
	defaultLength = 4
)

func (p *Parser) eofErr() error {
	return UnexpecedEofError{Row: p.t.row, Col: p.t.col}
}

type Parser struct {
	t *Tokenizer
}

// Convention: tokenizer always receives a tokenizer with a token available, excepting the Root
// program := constantDef* statement* ('loop:' statement*)?
func Parse(t *Tokenizer) (*song.Song, error) {
	p := &Parser{
		t: t,
	}
	s := &song.Song{
		Constants: map[string]song.Tablature{},
	}
	s.AddSyncedBlock()

	p.t.Next()
	for !p.t.EOF() {
		token := p.t.Get()
		switch token.Type {
		case ConstName:
			if err := p.constantDefNode(s); err != nil {
				return nil, err
			}
		case LoopTag:
			if err := p.loopNode(s); err != nil {
				return nil, err
			}
		case ChannelId, ChannelSync:
			if err := p.statementNode(s); err != nil {
				return nil, err
			}
		default:
			return nil, SyntaxError{t: token}
		}
		p.t.Next()
	}
	return s, nil
}

// constantDef := ID ':=' (instrumentDef | tablature+)
func (p *Parser) constantDefNode(s *song.Song) error {
	tok := p.t.Get()
	id := tok.getConstID()
	if _, ok := s.Constants[id]; ok {
		return RedefinitionError{tok}
	}
	if !p.t.Next() {
		return p.eofErr()
	}
	tok = p.t.Get()
	if tok.Type != Assign {
		return SyntaxError{tok}
	}
	if !p.t.Next() {
		return p.eofErr()
	}
	tok = p.t.Get()
	switch tok.Type {
	case OpenKey:
		inst, err := p.instrumentDefinitionNode(s)
		if err != nil {
			return err
		}
		s.Constants[id] = song.Tablature{{Instrument: &inst}}
	default:
		tabl, err := p.tablatureNode(s)
		if err != nil {
			return err
		}
		s.Constants[id] = tabl
	}
	return nil
}

// ('loop:' statement*)
func (p *Parser) loopNode(s *song.Song) error {
	return nil
}

// statement := channelFill | SYNC
func (p *Parser) statementNode(s *song.Song) error {
	return nil
}

// instrumentDef := '{' mapEntry* ('adsr:' adsrVector)? mapEntry* '}'
func (p *Parser) instrumentDefinitionNode(s *song.Song) (song.Instrument, error) {
	inst := song.Instrument{}
	if !p.t.Next() {
		return inst, p.eofErr()
	}
	definedAdsr, definedWave := false, false
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case AdsrVector:
			if definedAdsr {
				return inst, ParserError{tok, "defining ADSR envelope twice"}
			}
			definedAdsr = true
			attackLevel := float64(atoi(tok.Submatch[1])) / 100.0
			decayLevel := float64(atoi(tok.Submatch[3])) / 100.0
			inst.Envelope = []song.TimePoint{
				{Time: time.Duration(atoi(tok.Submatch[0])) * time.Millisecond, Val: attackLevel},
				{Time: time.Duration(atoi(tok.Submatch[2])) * time.Millisecond, Val: decayLevel},
				{Time: time.Duration(atoi(tok.Submatch[4])) * time.Millisecond, Val: decayLevel},
				{Time: time.Duration(atoi(tok.Submatch[5])) * time.Millisecond, Val: 0},
			}
		case MapEntry:
			switch strings.ToLower(tok.Submatch[0]) {
			case "adsr":
				return inst, ParserError{tok, "adsr should have a value like: 20->100, 50->80, 100, 120"}
			case "wave":
				if definedWave {
					return inst, ParserError{tok, "wave is defined twice"}
				}
				definedWave = true
				// todo: maybe validate wave values?
				inst.Wave = tok.Submatch[1]
			default:
				return inst, ParserError{tok, "only 'adsr' and 'wave' properties are allowed"}
			}
		case CloseKey:
			return inst, nil
		default:
			return inst, SyntaxError{tok}
		}
		p.t.Next()
	}
	return inst, nil
}

// panics as the regexp should have avoided filtering any unparseable number
func atoi(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		panic(fmt.Sprintf("Wrong number %q! This is a bug: %s", num, err.Error()))
	}
	return n
}

// constantDef := ID ':=' (instrumentDef | tablature+)
func (p *Parser) tablatureNode(s *song.Song) (song.Tablature, error) {
	return song.Tablature{}, nil
}
