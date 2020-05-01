package lang

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

type TokenType int

const (
	Channel TokenType = iota
	ChannelSendArrow
	String // any string, e.g. a tablature description
)

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

var tokens = regexp.MustCompile(`(@\w+)|(<-)|\S+|\|+`)

func (t *Tokenizer) Next() bool {
	for !t.EOF() {
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

func (t *Tokenizer) readMoreLines()  {
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

var channel = regexp.MustCompile(`^@(\w+)$`)

const arrow = "<-"

func parseToken(token []byte) Token {
	if ch := channel.FindSubmatch(token); ch != nil {
		return Token{Type: Channel, Content: token}
	}
	if string(token) == arrow {
		return Token{Type: ChannelSendArrow, Content: token}
	}
	// todo: verify the correct format of the string
	return Token{Type: String, Content: token}
}
