package lang

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestRedefinitionError(t *testing.T) {
	_, err := Parse(NewTokenizer(strings.NewReader(`
$foo := { wave: square }
$bar := abce
$foo := ffe
`)))
	require.IsTypef(t, RedefinitionError{}, err, "%#v", err)
	terr := err.(RedefinitionError)
	assert.Equal(t, 4, terr.t.Row)
	assert.Equal(t, 1, terr.t.Col)
}