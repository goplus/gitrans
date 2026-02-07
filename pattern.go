package gitrans

import (
	"path/filepath"
	"strings"
)

// -----------------------------------------------------------------------------

// matchResult defines outcomes of a match, no match, exclusion or inclusion.
type matchResult int

const (
	notMatched matchResult = iota
	matched
	exclude
)

// -----------------------------------------------------------------------------

const (
	excludePrefix  = "!"
	zeroToManyDirs = "**"
	patternDirSep  = "/"
)

type patternImpl struct {
	parts      []string
	zeroToMany int // index of ** if exists, -1 otherwise
	exclude    bool
}

// parsePattern parses a gitignore pattern string into the Pattern structure.
func parsePattern(p string) *patternImpl {
	res := &patternImpl{
		zeroToMany: -1,
	}

	if strings.HasPrefix(p, excludePrefix) {
		res.exclude = true
		p = p[1:]
	}

	res.parts = strings.Split(p, patternDirSep)
	for i, part := range res.parts {
		if part == zeroToManyDirs {
			res.zeroToMany = i
		}
	}
	return res
}

func parsePatterns(patterns []string) []*patternImpl {
	res := make([]*patternImpl, len(patterns))
	for i, p := range patterns {
		res[i] = parsePattern(p)
	}
	return res
}

func matchParts(parts, path []string) bool {
	for i, part := range parts {
		if match, err := filepath.Match(part, path[i]); err != nil || !match {
			return false
		}
	}
	return true
}

func (p *patternImpl) Match(path []string) matchResult {
	parts := p.parts
	if p.zeroToMany < 0 {
		if len(parts) != len(path) || !matchParts(parts, path) {
			return notMatched
		}
	} else {
		// Handle ** pattern matching
		if len(path) < len(parts)-1 || !matchParts(parts[:p.zeroToMany], path) {
			return notMatched
		}
		after := parts[p.zeroToMany+1:]
		if !matchParts(after, path[len(path)-len(after):]) {
			return notMatched
		}
	}
	if p.exclude {
		return exclude
	}
	return matched
}

func matchPattern(pattern []*patternImpl, name string) bool {
	path := strings.Split(name, patternDirSep)
	for _, p := range pattern {
		switch p.Match(path) {
		case matched:
			return true
		case exclude:
			return false
		}
	}
	return false
}

// -----------------------------------------------------------------------------
