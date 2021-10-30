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

type TokenType int

// ordered by order of precedence in the tokenization process (in case of multiple tokens match a string)
const (
	Comment TokenType = iota
	OpenInstrument
	SendArrow
	LoopTag
	OpenTuple
	CloseTuple
	CloseInstrument
	MapEntry
	Separator
	ConstDef
	ConstRef
	Assign
	ChannelId
	ChannelSync
	// Tablature stuff needs to go at the bottom, to not get confusion with other language grammar items
	Note
	Volume
	Silence
	Octave
	OctaveStep
	Number
	// NoMatch must be the last token
	NoMatch
)

func (t TokenType) String() string {
	switch t {
	case NoMatch:
		return "NoMatch"
	case LoopTag:
		return "LoopTag"
	case ConstDef:
		return "ConstDef"
	case ConstRef:
		return "ConstRef"
	case Assign:
		return "Assign"
	case OpenInstrument:
		return "OpenInstrument"
	case OpenTuple:
		return "OpenTuple"
	case CloseInstrument:
		return "CloseInstrument"
	case CloseTuple:
		return "CloseTuple"
	case MapEntry:
		return "MapEntry"
	case Separator:
		return "Separator"
	case ChannelSync:
		return "ChannelSync"
	case Comment:
		return "Comment"
	case Note:
		return "Note"
	case Volume:
		return "Volume"
	case Silence:
		return "Silence"
	case Octave:
		return "Octave"
	case OctaveStep:
		return "OctaveStep"
	case Number:
		return "Number"
	case ChannelId:
		return "ChannelId"
	case SendArrow:
		return "SendArrow"
	}
	return fmt.Sprintf("unknown: %d (probably a bug)", t)
}

var tokenDefs = map[TokenType]*regexp.Regexp{
	Comment:         regexp.MustCompile(`^;\.*$`),
	OpenInstrument:  regexp.MustCompile(`^(\w+)\s*\{$`),
	SendArrow:       regexp.MustCompile(`^<-$`),
	LoopTag:         regexp.MustCompile(`^[Ll][Oo][Oo][Pp]\s*:$`),
	OpenTuple:       regexp.MustCompile(`^\($`),
	CloseTuple:      regexp.MustCompile(`^\)(\d)+$`),
	CloseInstrument: regexp.MustCompile(`^}$`),
	MapEntry:        regexp.MustCompile(`^(\w+)\s*:\s*([^}#\n]+)$`),
	Separator:       regexp.MustCompile(`^\|+$`),
	ConstDef:        regexp.MustCompile(`^\$(\w+)\s*:=$`),
	ConstRef:        regexp.MustCompile(`^\$(\w+)$`),
	Assign:          regexp.MustCompile(`^:=$`),
	ChannelId:       regexp.MustCompile(`^@(\w+)$`),
	ChannelSync:     regexp.MustCompile(`^-{2,}$`),
	// Tablature stuff needs to go at the bottom, to not get confusion with other language grammar items
	Note:       regexp.MustCompile(`^([a-gA-G])([#+\-]?)(\d*)(\.*)$`),
	Volume:     regexp.MustCompile(`^[Vv](\d*)$`),
	Silence:    regexp.MustCompile(`^[Rr](\d*)$`),
	Octave:     regexp.MustCompile(`^[Oo](\d)$`),
	OctaveStep: regexp.MustCompile(`^(<|>)$`),
	Number:     regexp.MustCompile(`^(\d+)$`),
}

type Tokenizer struct {
	row       int
	col       int
	input     *bufio.Reader
	lineRest  string //line that is being currently parsed
	lastMatch string
	tokens    *regexp.Regexp
}

func NewTokenizer(input io.Reader, startRow int) *Tokenizer {
	return &Tokenizer{
		input:  bufio.NewReader(input),
		tokens: mergeAllTokens(),
		row: startRow,
	}
}

func mergeAllTokens() *regexp.Regexp {
	tokens := make([]TokenType, NoMatch)
	for i := 0; i < len(tokens); i++ {
		tokens[i] = TokenType(i)
	}
	sb := strings.Builder{}
	sb.WriteString("(")
	for _, r := range tokens {
		regex := tokenDefs[r].String()
		sb.WriteString(regex[:len(regex)-1]) //removing trailing $
		sb.WriteString(")|(")
	}
	sb.WriteString(`^\S+)`) // catching anything else as "unknown token"
	return regexp.MustCompile(sb.String())
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

// Get a token from a token type, if len(tokens) == 0, it searches across all the tokens
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

// if len(tokens) == 0, it searches across all the tokens
func (t *Tokenizer) parseToken(token string) Token {
	for tt := TokenType(0) ; tt < NoMatch ; tt++ {
		td := tokenDefs[tt]
		submatches := td.FindStringSubmatch(token)
		if submatches != nil {
			return Token{Type: tt, Content: token, Submatch: submatches[1:], Row: t.row, Col: t.col}
		}
	}
	return Token{Type: NoMatch, Content: token, Row: t.row, Col: t.col}
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
