package book

import (
	"regexp"
)

// submatchAsMap searches a string for regexp matches. Results are provided as
// a map whose keys are the regexp capture group name that triggered the
// submatch.
func submatchAsMap(s string, re *regexp.Regexp) map[string]string {
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
