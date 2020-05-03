package lang

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type TokenType int

const (
	Unknown TokenType = iota
	ChannelID
	ChannelSendArrow
	RampArrow
	OpenSection
	CloseSection
	CloseTuplet
	Separator
	Note
	Pause
	Octave
	IncOctave
	DecOctave
)

var tokenDefs = []struct {
	t TokenType
	r *regexp.Regexp
}{
	{t: ChannelID, r: regexp.MustCompile(`^@\w+$`)},
	{t: ChannelSendArrow, r: regexp.MustCompile(`^<-$`)},
	{t: RampArrow, r: regexp.MustCompile(`^->$`)},
	{t: OpenSection, r: regexp.MustCompile(`^\{$`)},
	{t: CloseTuplet, r: regexp.MustCompile(`^\}\d+`)},
	{t: CloseSection, r: regexp.MustCompile(`^\}$`)},
	{t: Note, r: regexp.MustCompile(`^[a-gA-G][#+-]?\d*\.*$`)},
	{t: Pause, r: regexp.MustCompile(`^[Rr]\d*$`)},
	{t: Octave, r: regexp.MustCompile(`^[Oo]\d$`)},
	{t: IncOctave, r: regexp.MustCompile(`^>$`)},
	{t: DecOctave, r: regexp.MustCompile(`^<$`)},
	{t: Separator, r: regexp.MustCompile(`^\|+$`)},
}

var tokens *regexp.Regexp

func init() {
	sb := strings.Builder{}
	sb.WriteString("(")
	for _, td := range tokenDefs {
		regex := td.r.String()
		sb.WriteString(regex[:len(regex)-1]) //removing trailing $
		sb.WriteString(")|(")
	}
	sb.WriteString(`^\S+)`) // catching anything else as "unknown token"
	tokens = regexp.MustCompile(sb.String())
}

type Tokenizer struct {
	Row       int
	Col       int
	input     *bufio.Reader
	lineRest  []byte //line that is being currently parsed
	lastMatch []byte
}

func NewTokenizer(input io.Reader) *Tokenizer {
	return &Tokenizer{
		input: bufio.NewReader(input),
	}
}

func (t *Tokenizer) Next() bool {
	for !t.EOF() {
		// trimming leading spaces
		i := 0
		for i < len(t.lineRest) && (t.lineRest[i] == ' ' || t.lineRest[i] == '\t') {
			i++
		}
		t.Col += i
		t.lineRest = t.lineRest[i:]
		idx := tokens.FindIndex(t.lineRest)
		if idx != nil {
			t.lastMatch = t.lineRest[idx[0]:idx[1]]
			t.lineRest = t.lineRest[idx[1]:]
			t.Col += idx[0]
			return true
		}
		t.readMoreLines()
	}
	return false
}

func (t *Tokenizer) readMoreLines() {
	var err error
	t.lastMatch = nil
	t.lineRest, err = t.input.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			t.input = nil
			t.lineRest = nil
			return
		}
		panic(fmt.Errorf("can't read next line: %w", err))
	}
	t.Col = 1
	t.Row++
}

func (t *Tokenizer) EOF() bool {
	return len(t.lineRest) == 0 && t.input == nil
}

func (t *Tokenizer) Get() Token {
	return parseToken(t.lastMatch)
}

type Token struct {
	Type    TokenType
	Content []byte
}

func parseToken(token []byte) Token {
	for _, td := range tokenDefs {
		if td.r.Match(token) {
			return Token{Type: td.t, Content: token}
		}
	}
	return Token{Type: Unknown, Content: token}
}
