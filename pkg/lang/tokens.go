package lang

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/mariomac/msxmml/pkg/song/note"
)

type TokenType string

const (
	AnyString       TokenType = "AnyString"
	LoopTag         TokenType = "LoopTag"
	ConstDef        TokenType = "ConstDef"
	ConstRef        TokenType = "ConstRef"
	Assign          TokenType = "Assign"
	OpenInstrument  TokenType = "OpenInstrument"
	OpenTuple       TokenType = "OpenTuple"
	CloseInstrument TokenType = "CloseInstrument"
	CloseTuple      TokenType = "CloseTuple"
	MapEntry        TokenType = "MapEntry"
	Separator       TokenType = "Separator"
	ChannelSync     TokenType = "ChannelSync"
	Comment         TokenType = "Comment"
	Note            TokenType = "Note"
	Volume          TokenType = "Volume"
	Silence         TokenType = "Silence"
	Octave          TokenType = "Octave"
	OctaveStep      TokenType = "OctaveStep"
	Number          TokenType = "Number"
	ChannelId       TokenType = "ChannelId"
	SendArrow       TokenType = "SendArrow"
)

var tokenDefs = []struct {
	t TokenType
	r *regexp.Regexp
}{
	{t: Comment, r: regexp.MustCompile(`^;\.*$`)},
	{t: OpenInstrument, r: regexp.MustCompile(`^(\w+)\s*\{$`)},
	{t: SendArrow, r: regexp.MustCompile(`^<-$`)},
	{t: LoopTag, r: regexp.MustCompile(`^[Ll][Oo][Oo][Pp]\s*:$`)},
	{t: OpenTuple, r: regexp.MustCompile(`^\($`)},
	{t: CloseTuple, r: regexp.MustCompile(`^\)(\d)+$`)},
	{t: CloseInstrument, r: regexp.MustCompile(`^}$`)},
	{t: MapEntry, r: regexp.MustCompile(`^(\w+)\s*:\s*([^}#\n]+)$`)},
	{t: Separator, r: regexp.MustCompile(`^\|+$`)},
	{t: ConstDef, r: regexp.MustCompile(`^\$(\w+)\s*:=$`)},
	{t: ConstRef, r: regexp.MustCompile(`^\$(\w+)$`)},
	{t: Assign, r: regexp.MustCompile(`^:=$`)},
	{t: ChannelId, r: regexp.MustCompile(`^@(\w+)$`)},
	{t: ChannelSync, r: regexp.MustCompile(`^-{2,}$`)},
	// Tablature stuff needs to go at the bottom, to not get confusion with other language grammar items
	{t: Note, r: regexp.MustCompile(`^([a-gA-G])([#+\-]?)(\d*)(\.*)$`)},
	{t: Volume, r: regexp.MustCompile(`^[Vv](\d*)$`)},
	{t: Silence, r: regexp.MustCompile(`^[Rr](\d*)$`)},
	{t: Octave, r: regexp.MustCompile(`^[Oo](\d)$`)},
	{t: OctaveStep, r: regexp.MustCompile(`^(<|>)$`)},
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

// todo: ignore comments
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
	Type TokenType
	// TODO: replace content[0] invocations by typesafe functions
	Content string
	// TODO: replace inline indexing by typesafe functions
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

func (f *Token) assertType(expected TokenType) {
	if f.Type != expected {
		panic(fmt.Sprintf("BUG detected. Expected type: %s. Got: %s", expected, f.Type))
	}
}

func (f *Token) getConstDefId() string {
	f.assertType(ConstDef)
	return f.Submatch[0]
}

func (f *Token) getConstRefId() string {
	f.assertType(ConstRef)
	return f.Submatch[0]
}

func (f *Token) getTupletNumber() int {
	f.assertType(CloseTuple)
	return mustAtoi(f.Submatch[0])
}

func mustAtoi(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		panic(fmt.Sprintf("BUG detected. Expected number, got %q", num))
	}
	return n
}

func (f *Token) getOctaveStep() int {
	f.assertType(OctaveStep)
	switch f.Content[0] {
	case '<':
		return -1
	case '>':
		return +1
	default:
		panic(fmt.Sprintf("BUG detected. Invalid octave step %q", f.Content))
	}
}

var pitches = [8]note.Pitch{note.A, note.B, note.C, note.D, note.E, note.F, note.G}

// A note should come represented by an array where
// 0: pitch - 1: halftone - 2: length - 3: dots
// todo: return t, error if a given note can't be sharp or flat
func (f *Token) getNote() (note.Note, error) {
	f.assertType(Note)

	var pitch note.Pitch
	c := f.Submatch[0][0]
	if c >= 'A' && c <= 'Z' {
		pitch = pitches[c-'A']
	} else if c >= 'a' && c <= 'z' {
		pitch = pitches[c-'a']
	} else {
		panic(fmt.Sprintf("BUG detected. Pitch can't be '%c'", c))
	}

	n := note.Note{
		Pitch:    pitch,
		Length:   defaultLength,
		Halftone: note.NoHalftone,
		Dots:     len(f.Submatch[3]),
	}
	// get halftone
	if len(f.Submatch[1]) > 0 {
		switch f.Submatch[1][0] {
		case '+', '#':
			n.Halftone = note.Sharp
		case '-':
			n.Halftone = note.Flat
		default:
			panic(fmt.Sprintf("BUG detected. Wrong halftone %q", f.Submatch[1]))
		}
	}

	// get Length
	if len(f.Submatch[2]) > 0 {
		l, err := strconv.Atoi(f.Submatch[2])
		if err != nil {
			panic(fmt.Sprintf("BUG detected. Wrong length for note: %#v. Err: %s",
				f, err.Error()))
		}
		if l < minLength || l > maxLength {
			return n, fmt.Errorf(
				"wrong note length: %d. Must be in range %d to %d", l, minLength, maxLength)
		}
		n.Length = l
	}
	return n, nil
}

func (token *Token) getOctave() int {
	token.assertType(Octave)
	return mustAtoi(token.Submatch[0])
}

func (token *Token) getVolume() int {
	token.assertType(Volume)
	return mustAtoi(token.Submatch[0])
}

func (token *Token) getSilence() note.Note {
	token.assertType(Silence)
	n := note.Note{Pitch: note.Silence}
	if len(token.Submatch[0]) == 0 {
		n.Length = defaultLength
		return n
	}
	n.Length = mustAtoi(token.Submatch[0])
	return n
}

func (tok *Token) getInstrumentClass() string {
	tok.assertType(OpenInstrument)
	return tok.Submatch[0]
}

func (tok *Token) getMapKeyValue() (string, string) {
	tok.assertType(MapEntry)
	return tok.Submatch[0], tok.Submatch[1]
}

func (t *Token) getChannelId() string {
	t.assertType(ChannelId)
	return t.Submatch[0]
}
