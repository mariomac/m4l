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
$bar := 1234
$foo := 3321
`)))
	require.IsType(t, RedefinitionError{}, err)
	terr := err.(RedefinitionError)
	assert.Equal(t, 4, terr.t.Row)
	assert.Equal(t, 1, terr.t.Col)
}