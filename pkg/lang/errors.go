package lang

import "fmt"

func errHeader(t Token) string {
	return fmt.Sprintf("%d:%d - Unexpected EOF", t.Row, t.Col)
}

type SyntaxError struct {
	t Token
}

func (p SyntaxError) Error() string {
	return errHeader(p.t) + fmt.Sprintf("Unexpected %q", p.t.Content)
}

type UnexpecedEofError struct {
	Row int
	Col int
}

func (p UnexpecedEofError) Error() string {
	return fmt.Sprintf("%d:%d - Unexpected EOF", p.Row, p.Col)
}

type RedefinitionError struct {
	t Token
}

func (r RedefinitionError) Error() string {
	return errHeader(r.t) + fmt.Sprintf("can't redefine: %v", r.t.Content)
}

type ParserError struct {
	t   Token
	msg string
}

func (p ParserError) Error() string {
	return errHeader(p.t) + p.msg
}
