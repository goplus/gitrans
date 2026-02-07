package gitrans

import (
	"testing"
)

// -----------------------------------------------------------------------------

func TestPattern(t *testing.T) {
	type patternCase struct {
		pattern  []string
		path     string
		expected bool
	}
	cases := []patternCase{
		{[]string{"**"}, "a/b/c", true},
		{[]string{"b", "**/c"}, "a/b/c", true},
		{[]string{"*/*/c"}, "a/b/c", true},
		{[]string{"b/**", "**/b"}, "a/b/c", false},
		{[]string{"!*/*/c", "**"}, "a/b/c", false},
		{[]string{"a", "c"}, "a/b/c", false},
	}
	for _, c := range cases {
		p := parsePatterns(c.pattern)
		if matchPattern(p, c.path) != c.expected {
			t.Errorf("pattern %v on path %q: expected %v", c.pattern, c.path, c.expected)
		}
	}
}

// -----------------------------------------------------------------------------
