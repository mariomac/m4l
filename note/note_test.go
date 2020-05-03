package note

//func TestParseChannel(t *testing.T) {
//	_, err := Parse([]byte("a16b32co3a<b#3c-"))
//	require.NoError(t, err)
//}
//
//func TestParseChannel_WrongStrings(t *testing.T) {
//	testCases := []struct {
//		desc string
//		input string
//		errPrefix string
//	} {
//		{ "length after halftone", "O5cde#88", "at position 4" }, // it's position 6 but the note starts in position 4
//		{ "octave without number", "cBaO", "at position 3"},
//		{ "octave below minimum", "<<<<", "at position 3"},
//		{ "octave over maximum", "o6>>>", "at position 4"},
//		{ "wrong octave", "o99", "at position 0"},
//		{ "wrong octave 0", "o0", "at position 0"},
//		{ "two halftones", "A<<b--", "at position 5"},
//		{ "wrong note length 0", "o5a0", "at position 2"},
//		{ "wrong note length 1000", "O5a1000", "at position 2"},
//		{ "invalid characters at the end", "abcde!!!", "at position 5"},
//		{ "invalid notes", "abcdedgo2>ab#4abcdeebbfghwa38", "unknown char"},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.desc, func(t *testing.T) {
//			_, err := Parse([]byte(tc.input))
//			require.Error(t, err)
//			require.Truef(t, strings.HasPrefix(err.Error(), tc.errPrefix),
//				"expected prefix %q in error message %q", tc.errPrefix, err.Error())
//		})
//	}
//}
