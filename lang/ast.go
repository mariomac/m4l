package lang
//
//import (
//	"fmt"
//	"regexp"
//
//	"github.com/mariomac/msxmml/solfa"
//
//	"github.com/mariomac/msxmml/song"
//)
//
//type Parser struct {
//	t *Tokenizer
//}
//
//func NewParser(t *Tokenizer) *Parser {
//	return &Parser{
//		t:t,
//	}
//}
//
//// Convention: tokenizer always receives a tokenizer with a token available, excepting the Root
//
//// song: channel+
//func (p *Parser) Parse() (*song.Song, error) {
//	s := &song.Song{Channels: map[string]song.Channel{}}
//	p.t.Next()
//	for !p.t.EOF() {
//		token := p.t.Get()
//		switch token.Type {
//		case ChannelID:
//			ch, err := p.channelNode()
//			if err != nil {
//				return nil, err
//			}
//			s.Channels[ch.Name] = ch
//		default:
//			return nil, &SyntaxError{t: token}
//		}
//	}
//	return s, nil
//}
//
//func (p *Parser) channelNode() (song.Channel, error) {
//	last := p.t.Get()
//	c := song.Channel{Name: string(last.Content[1:])}
//
//	if !p.t.Next() {
//		return c, unexpectedError(t, "channel information", []byte("end of input"))
//	}
//	last = p.t.Get()
//	if last.t.Type != ChannelSendArrow {
//		return c, unexpectedError(t, "an arrow '<-'", last.t.Content)
//	}
//	if !p.t.Next() {
//		return c, unexpectedError(t, "channel information", []byte("end of input"))
//	}
//	tabs, err := TablatureNode(t)
//	if err != nil {
//		return c, err
//	}
//	// TODO: keep last note/global config for each channel number
//	// so octave and other data stays between successive commands
//	c.Notes, err = solfa.Parse(tabs)
//	if err != nil {
//		return c, &ParserError{t: t, cause: fmt.Errorf("problem with channel tablature: %w", err)}
//	}
//	return c, nil
//}
//
//var tabRegex = regexp.MustCompile(`^(([a-zA-Z][+\-#]?\d*)|[<>]|\||(\{[^\{}]*}\d+))*$`)
//
//func  (p *Parser) TablatureNode() ([]byte, error) {
//	var tablature []byte
//	tok := p.t.Get()
//	if tok.Type != 0 && !tabRegex.Match(tok.Content) {
//		return nil, unexpectedError(t, "a music tablature", tok.Content)
//	}
//	tablature = append(tablature, tok.Contenp.t...)
//	for p.t.Next() {
//		tok := p.t.Get()
//		if tok.Type != 0 && !tabRegex.Match(tok.Content) {
//			return tablature, nil
//		}
//		tablature = append(tablature, tok.Contenp.t...)
//	}
//
//	return tablature, nil
//}
//
//type SyntaxError struct {
//	t     Token
//}
//
//func (p *SyntaxError) Error() string {
//	return fmt.Sprintf("%d:%d - Unexpected %q", p.t.Row, p.t.Col, string(p.t.Content))
//}
