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
	Silence
	Octave
	IncOctave
	DecOctave
)

var tokenDefs = []struct {
	t TokenType
	r *regexp.Regexp
}{
	{t: ChannelID, r: regexp.MustCompile(`^@(\w+)$`)},
	{t: ChannelSendArrow, r: regexp.MustCompile(`^<-$`)},
	{t: RampArrow, r: regexp.MustCompile(`^->$`)},
	{t: OpenSection, r: regexp.MustCompile(`^\{$`)},
	{t: CloseTuplet, r: regexp.MustCompile(`^}(\d+)$`)},
	{t: CloseSection, r: regexp.MustCompile(`^}$`)},
	{t: Note, r: regexp.MustCompile(`^([a-gA-G])([#+-]?)(\d*)(\.*)$`)},
	{t: Silence, r: regexp.MustCompile(`^[Rr](\d*)$`)},
	{t: Octave, r: regexp.MustCompile(`^[Oo](\d)$`)},
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
	row       int
	col       int
	input     *bufio.Reader
	lineRest  string //line that is being currently parsed
	lastMatch string
}

func NewTokenizer(input io.Reader) *Tokenizer {
	return &Tokenizer{
		input: bufio.NewReader(input),
	}
}

func (t *Tokenizer) Next() bool {
	t.col += len(t.lastMatch)
	for !t.EOF() {
		// trimming leading spaces
		i := 0
		for i < len(t.lineRest) && (t.lineRest[i] == ' ' || t.lineRest[i] == '\t') {
			i++
		}
		t.col += i
		t.lineRest = t.lineRest[i:]
		idx := tokens.FindStringIndex(t.lineRest)
		if idx != nil {
			t.lastMatch = t.lineRest[idx[0]:idx[1]]
			t.lineRest = t.lineRest[idx[1]:]
			return true
		}
		t.readMoreLines()
	}
	return false
}

func (t *Tokenizer) readMoreLines() {
	var err error
	t.lastMatch = ""
	t.lineRest, err = t.input.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			t.input = nil
			t.lineRest = ""
			return
		}
		panic(fmt.Errorf("can't read next line: %w", err))
	}
	t.col = 1
	t.row++
}

func (t *Tokenizer) EOF() bool {
	return len(t.lineRest) == 0 && t.input == nil
}

func (t *Tokenizer) Get() Token {
	return t.parseToken(t.lastMatch)
}

type Token struct {
	Type     TokenType
	Content  string
	Submatch []string
	Row, Col int
}

func (t *Tokenizer) parseToken(token string) Token {
	for _, td := range tokenDefs {
		submatches := td.r.FindStringSubmatch(token)
		if submatches != nil {
			return Token{Type: td.t, Content: token, Submatch: submatches[1:], Row: t.row, Col: t.col}
		}
	}
	return Token{Type: Unknown, Content: token, Row: t.row, Col: t.col}
}
