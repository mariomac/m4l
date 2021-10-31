package lang

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedefinitionError(t *testing.T) {
	_, err := Parse(strings.NewReader(`
$foo := pgs { wave: square }
$bar := abce
$foo := ffe
`))
	require.IsTypef(t, RedefinitionError{}, err, "%#v", err)
	terr := err.(RedefinitionError)
	assert.Equal(t, 4, terr.t.Row)
	assert.Equal(t, 1, terr.t.Col)
}

func TestSyntaxError(t *testing.T) {
	_, err := Parse(strings.NewReader(`
$foo := ( wave: square )
`))
	require.IsTypef(t, SyntaxError{}, err, "%#v", err)
	terr := err.(SyntaxError)
	assert.Equal(t, 2, terr.t.Row)
	assert.Equal(t, 11, terr.t.Col)
}

func TestConstantIntoConstant(t *testing.T) {
	_, err := Parse(strings.NewReader(`
$foo := abc
$bar := a $foo b 
`))
	assert.Error(t, err)
	assert.IsType(t, ParserError{}, err)
	assert.Equal(t, 3, err.(ParserError).t.Row)
	assert.Equal(t, 11, err.(ParserError).t.Col)
}
