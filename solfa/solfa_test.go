package solfa

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseChannel(t *testing.T) {
	_, err := ParseChannel([]byte("a16b32co3a<b3#c-"))
	require.NoError(t, err)
}

func TestParseChannel_WrongStrings(t *testing.T) {
	testCases := []struct {
		desc string
		input string
		errPrefix string
	} {
		{ "length after halftone", "O5cde+88", "at position 6" },
		{ "octave without number", "cBaO", "at position 3"},
		{ "octave below minimum", "<<<<", "at position 3"},
		{ "octave over maximum", "o6>>>", "at position 4"},
		{ "wrong octave", "o99", "at position 0"},
		{ "wrong octave 0", "o0", "at position 0"},
		{ "two halftones", "A<<b--", "at position 5"},
		{ "wrong note length 0", "o5a0", "at position 2"},
		{ "wrong note length 1000", "O5a1000", "at position 2"},

	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := ParseChannel([]byte(tc.input))
			require.Error(t, err)
			require.Truef(t, strings.HasPrefix(err.Error(), tc.errPrefix),
				"expected prefix %q in error message %q", tc.errPrefix, err.Error())
		})
	}
}
