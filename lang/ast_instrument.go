package lang

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mariomac/msxmml/song"
)

// instrumentDefinition: 'wave:' \w+ \n 'adsr:' num->num, num->num, num, num
func (p *Parser) instrumentDefinitionNode(c *song.Channel) error {
	if !p.t.Next() {
		return p.eofErr()
	}
	definedAdsr, definedWave := false, false
	for !p.t.EOF() {
		tok := p.t.Get()
		switch tok.Type {
		case ADSRVector:
			if definedAdsr {
				return ParserError{tok, "defining ADSR envelope twice"}
			}
			definedAdsr = true
			attackLevel := float64(atoi(tok.Submatch[1])) / 100.0
			decayLevel := float64(atoi(tok.Submatch[3])) / 100.0
			c.Instrument.Envelope = []song.TimePoint{
				{Time: time.Duration(atoi(tok.Submatch[0])) * time.Millisecond, Val: attackLevel},
				{Time: time.Duration(atoi(tok.Submatch[2])) * time.Millisecond, Val: decayLevel},
				{Time: time.Duration(atoi(tok.Submatch[4])) * time.Millisecond, Val: decayLevel},
				{Time: time.Duration(atoi(tok.Submatch[5])) * time.Millisecond, Val: 0},
			}
		case MapEntry:
			switch strings.ToLower(tok.Submatch[0]) {
			case "adsr":
				return ParserError{tok, "adsr should have a value like: 20->100, 50->80, 100, 120"}
			case "wave":
				if definedWave {
					return ParserError{tok, "wave is defined twice"}
				}
				definedWave = true
				// todo: maybe validate wave values?
				c.Instrument.Wave = tok.Submatch[1]
			default:
				return ParserError{tok, "only 'adsr' and 'wave' properties are allowed"}
			}
		case CloseSection:
			p.t.Next()
			return nil
		default:
			return SyntaxError{tok}
		}
		p.t.Next()
	}

	return nil
}

// panics as the regexp should have avoided filtering any unparseable number
func atoi(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		panic(fmt.Sprintf("Wrong number %q! This is a bug: %s", num, err.Error()))
	}
	return n
}
