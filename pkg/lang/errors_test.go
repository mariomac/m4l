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
