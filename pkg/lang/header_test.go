package lang

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWithHeader(t *testing.T) {
	s, err := Parse(strings.NewReader(`
; some comments here
psg.hz 60
tempo 120     ; this should be ignored

; program starts here
@ch1 <- abc
`))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"psg.hz": "60",
		"tempo":  "120",
	}, s.Properties)
	require.Len(t, s.Blocks, 1)
	require.Len(t, s.Blocks[0].Channels["ch1"].Items, 3)
}

func TestParseWithHeader_Tokenizer_Position_Ok(t *testing.T) {
	_, err := Parse(strings.NewReader(`
; some comments here
psg.hz 60
tempo 120     ; this should be ignored

; should show an error here
@ch1 <- tracatraca
`))
	assert.IsType(t, SyntaxError{}, err)
	assert.Equal(t, err.(SyntaxError).t.Col, 9)
	assert.Equal(t, err.(SyntaxError).t.Row, 7)
}
