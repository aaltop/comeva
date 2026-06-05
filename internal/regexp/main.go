package regexp

import (
	baseRegexp "regexp"
)

type SubMatch struct {
	// Match is the matched string. Empty if not match (or the matched empty).
	Match string
	// Cols contains the start and end indices (columns) of the match. If the
	// indices are minus one (-1), no match was found.
	Cols [2]int
}

// GetSubMatches returns the found submatches. The keys of the map correspond
// to the names of the parenthesised subexpressions, as well as the empty string
// ("") for the whole found expression.
func GetSubMatches(re *baseRegexp.Regexp, s string) (matches map[string]SubMatch) {
	var stringMatches = re.FindStringSubmatch(s)
	var indexMatches = re.FindStringSubmatchIndex(s)

	var ind int
	matches = make(map[string]SubMatch)
	for _, name := range re.SubexpNames() {
		ind = re.SubexpIndex(name)
		if name == "" {
			ind = 0
		}
		matches[name] = SubMatch{
			Match: stringMatches[ind],
			Cols:  [2]int{indexMatches[2*ind], indexMatches[2*ind+1]},
		}
	}
	return
}
