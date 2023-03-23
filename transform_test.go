package diskv

import (
	"testing"
)

func TestTransform(t *testing.T) {
	var (
		prefixFn    = BlockTransform(3, 3, false)
		prefixCases = map[string][]string{
			"0":     {},
			"diskv": {"dis"},
			"QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH":              {"Qmb", "FMk", "e1K"},
			"bafybeifx7yeb55armcsxwwitkymga5xf53dxiarykms3ygqic223w5sk3m": {"baf", "ybe", "ifx"},
		}

		suffixFn    = BlockTransform(3, 3, true)
		suffixCases = map[string][]string{
			"0":     {},
			"diskv": {"skv"},
			"QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH":              {"Gu6", "x1A", "wQH"},
			"bafybeifx7yeb55armcsxwwitkymga5xf53dxiarykms3ygqic223w5sk3m": {"223", "w5s", "k3m"},
		}
	)

	for s, expected := range prefixCases {
		if value := prefixFn(s); !cmpStrings(expected, value) {
			t.Fatalf("prefix transform %s expected: %v, got: %v", s, expected, value)
		}
	}

	for s, expected := range suffixCases {
		if value := suffixFn(s); !cmpStrings(expected, value) {
			t.Fatalf("suffix transform %s expected: %v, got: %v", s, expected, value)
		}
	}
}
