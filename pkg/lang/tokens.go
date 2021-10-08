package lang

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type TokenType string

const (
	AnyString TokenType = "AnyString"
	LoopTag TokenType = "LoopTag"
	ConstName TokenType = "ConstName"
	Assign TokenType = "Assign"
	OpenKey TokenType = "OpenKey"
	CloseKey TokenType = "CloseKey"
	MapEntry TokenType = "MapEntry"
	AdsrVector TokenType = "AdsrVector"
	Separator TokenType = "Separator"
	ChannelSync TokenType = "ChannelSync"
	Comment TokenType = "Comment"
	Note TokenType = "Note"
	Silence TokenType = "Silence"
	Octave TokenType = "Octave"
	IncOctave TokenType = "OctaveStep"
	DecOctave TokenType = "DecOctave"
	Number TokenType = "Number"
	ChannelId TokenType = "ChannelId"
	SendArrow TokenType = "SendArrow"
)

var tokenDefs = []struct {
	t TokenType
	r *regexp.Regexp
}{
	{t: Comment, r: regexp.MustCompile(`^#\.*$`) },
	{t: SendArrow, r: regexp.MustCompile(`^<-$`)},
	{t: LoopTag, r: regexp.MustCompile(`^[Ll][Oo][Oo][Pp]\s*:$`)},
	{t: OpenKey, r: regexp.MustCompile(`^\{$`)},
	{t: CloseKey, r: regexp.MustCompile(`^}$`)},
	{t: AdsrVector, r: regexp.MustCompile(`^[Aa][Dd][Ss][Rr]\s*:\s*(\d+)\s*->\s*(\d+)\s*,\s*(\d+)\s*\->\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)$`)},
	{t: MapEntry, r: regexp.MustCompile(`^(\w+)\s*:\s*(\w+)$`)},
	{t: Separator, r: regexp.MustCompile(`^\|+$`)},
	{t: ConstName, r: regexp.MustCompile(`^\$(\w+)$`)},
	{t: Assign, r: regexp.MustCompile(`^:=$`)},
	{t: ChannelId, r: regexp.MustCompile(`^@(\w+)$`)},
	{t: ChannelSync, r: regexp.MustCompile(`^-{2,}$`)},
	// Tablature stuff needs to go at the bottom, to not get confusion with other language grammar items
	{t: Note, r: regexp.MustCompile(`^([a-gA-G])([#+\-]?)(\d*)(\.*)$`)},
	{t: Silence, r: regexp.MustCompile(`^[Rr](\d*)$`)},
	{t: Octave, r: regexp.MustCompile(`^[Oo](\d)$`)},
	{t: IncOctave, r: regexp.MustCompile(`^>$`)},
	{t: DecOctave, r: regexp.MustCompile(`^<$`)},
	{t: Number, r: regexp.MustCompile(`^(\d+)$`)},
}

type Tokenizer struct {
	row       int
	col       int
	input     *bufio.Reader
	lineRest  string //line that is being currently parsed
	lastMatch string
	tokens    *regexp.Regexp
}

func NewTokenizer(input io.Reader) *Tokenizer {
	sb := strings.Builder{}
	sb.WriteString("(")
	for _, r := range tokenDefs {
		regex := r.r.String()
		sb.WriteString(regex[:len(regex)-1]) //removing trailing $
		sb.WriteString(")|(")
	}
	sb.WriteString(`^\S+)`) // catching anything else as "unknown token"

	return &Tokenizer{
		input:  bufio.NewReader(input),
		tokens: regexp.MustCompile(sb.String()),
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
		idx := t.tokens.FindStringIndex(t.lineRest)
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
	return Token{Type: AnyString, Content: token, Row: t.row, Col: t.col}
}

func (f *Token) getConstID() string {
	f.assertType(ConstName)
	return f.Submatch[0]
}

func (f *Token) assertType(expected TokenType) {
	if f.Type != expected {
		panic(fmt.Sprintf("expected type: %s. Got: %s", expected, f.Type))
	}
}
