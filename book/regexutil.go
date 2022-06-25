package book

import (
	"bufio"
	"io"
	"regexp"
)

// reFindStringSubmatchAsMap searches a string for regexp matches. Results are
// provided as a map whose keys are the regexp capture group name that
// triggered the submatch.
func reFindStringSubmatchAsMap(s string, re *regexp.Regexp) map[string]string {
	matches, names := re.FindStringSubmatch(s), re.SubexpNames()
	if matches != nil {
		found := make(map[string]string, len(matches)-1)
		for i := range matches {
			if i > 0 {
				found[names[i]] = matches[i]
			}
		}

		return found
	}

	return nil
}

// reFindReaderSubmatchAsMap searches a io.Reader for regexp matches. Results are
// provided as a map whose keys are the regexp capture group name that
// triggered the submatch.
func reFindReaderSubmatchAsMap(r io.Reader, re *regexp.Regexp) (matches []map[string]string) {
	splitRe := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		loc, names := re.FindSubmatchIndex(data), re.SubexpNames()

		if hasMatches(loc) {
			found := make(map[string]string, len(names)-1)
			for i := range names {
				if i > 0 {
					found[names[i]] = string(data[loc[2*i]:loc[2*i+1]])
				}
			}
			matches = append(matches, found)
			return loc[1] + 1, data[loc[0]:loc[1]], nil
		}

		if atEOF {
			return 0, nil, io.EOF
		}

		return 0, nil, nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(splitRe)
	for scanner.Scan() {
	}

	return
}

func hasMatches(loc []int) bool {
	if len(loc) == 0 {
		return false
	}

	for _, i := range loc {
		if i < 0 {
			return false
		}
	}

	return true
}
