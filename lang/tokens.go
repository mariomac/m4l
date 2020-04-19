package lang

import (
	"fmt"
	"regexp"
	"strconv"
)

type TokenType int

const (
	Channel TokenType = iota
	ChannelSendArrow
	String // any string, e.g. a tablature description
)

type Tokenizer struct {
	rest      []byte // subslice to non-read input
	lastMatch []byte
}

func NewTokenizer(input []byte) *Tokenizer {
	return &Tokenizer{
		rest: input,
	}
}

var tokens = regexp.MustCompile(`\S+`)

func (t *Tokenizer) Next() bool {
	if len(t.rest) == 0 {
		return false
	}
	idx := tokens.FindIndex(t.rest)
	if idx == nil {
		t.lastMatch = nil
		t.rest = nil // EOF
		return false
	}
	t.lastMatch = t.rest[idx[0]:idx[1]]
	t.rest = t.rest[idx[1]:]
	return true
}

func (t *Tokenizer) Get() Token {
	return parseToken(t.lastMatch)
}

type Token struct {
	Type    TokenType
	Content []byte
	Channel ChannelToken // if Type == Channel, this has metadata about the channel
}

type ChannelToken struct {
	Number int
}

var channel = regexp.MustCompile(`^@(\d+)$`)

const arrow = "<-"

func parseToken(token []byte) Token {
	if ch := channel.FindSubmatch(token); ch != nil {
		num, err := strconv.Atoi(string(ch[1]))
		if err != nil {
			panic(fmt.Errorf("THIS IS BUG!!! expecting a channel number: %w", err))
		}
		return Token{Type: Channel, Content: token, Channel: ChannelToken{Number: num}}
	}
	if string(token) == arrow {
		return Token{Type: ChannelSendArrow, Content: token}
	}
	// todo: verify the correct format of the string
	return Token{Type: String, Content: token}
}
