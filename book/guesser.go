package book

import (
	"regexp"
)

const (
	// seriesKey lists usual keywords that introduce a specification of a
	// series index.
	reSeriesIndex = `(?i:Series |Volume |Vol |Vol. |Part |n°|#|)(?P<SeriesIndex>\d+)`
)

var (
	// pathGuesser is a collection of regexp to extract information from a
	// Book's filename.
	// Order is important as only first match is considered, so it needs to be
	// defined from the more specific to the more general capture logic.
	pathGuesser = reGuesser{
		// parent/folder/<Authors> - [<Series> <SeriesIndex>] - <Title> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s-\s\[(?P<Series>.+)\s(?P<SeriesIndex>\d+)\]\s-\s(?P<Title>.+?)\s\[(?P<Language>.+)\]\.(?:.+)$`),
		// parent/folder/<Authors> - <Title> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s-\s(?P<Title>.+?)\s\[(?P<Language>.+)\]\.(?:.+)$`),
	}

	// seriesGuesser is a collection of regexp to extract series information
	// from a Book's title or subtitle.
	// Order is important as only first match is considered, so it needs to be
	// defined from the more specific to the more general capture logic.
	seriesGuesser = reGuesser{
		// <ShortTitle> (<Series> n°<SeriesIndex>)
		regexp.MustCompile(`^(?P<ShortTitle>.+)\s\p{Ps}(?P<Series>.+?)\s` + reSeriesIndex + `\p{Pe}$`),
		// <ShortTitle> - <Series> n°<SeriesIndex>
		regexp.MustCompile(`^(?P<ShortTitle>.+?)\s\p{Pd}\s(?P<Series>.+?)\s` + reSeriesIndex + `$`),
		// <Series> n°<SeriesIndex>
		regexp.MustCompile(`^(?P<Series>.+?)\s` + reSeriesIndex + `$`),
		// Book <SeriesIndex> of <Series>
		regexp.MustCompile(`^Book\s(?P<SeriesIndex>\d+)\sof\s(?P<Series>.+)$`),
	}

	// titleGuesser is a collection of regexp that pre-processes bad-formatted Titles
	// Order is important as only first match is considered, so it needs to be
	// defined from the more specific to the more general capture logic.
	titleGuesser = reGuesser{
		// <ShortTitle> / <SubTitle>
		regexp.MustCompile(`^(?P<ShortTitle>.+)\s/\s(?P<SubTitle>.+)$`),
	}
)

// reGuesser represents a guesser based on regexp.
type reGuesser []*regexp.Regexp

// GuessFrom extracts fields from a string based on guesser's regexp
// collection. Regexps are run in the guesser declaration order and first match
// is returned.
func (g reGuesser) GuessFrom(s string) map[string]string {
	for _, re := range g {
		matches, names := re.FindStringSubmatch(s), re.SubexpNames()
		if matches != nil {
			found := make(map[string]string, len(matches)-1)
			for i := range matches {
				if i > 0 {
					found[names[i]] = matches[i]
				}
			}

            Debug.Printf("guessed information: '%+v'", found)
			return found
		}
	}

	return nil
}

// GuessFromPath tries to guess some of Book's fields based on Book's path.
func (b *Book) GuessFromPath() {
	Debug.Printf("Guess information from book's path '%s'", b.Path)
    guessed := pathGuesser.GuessFrom(b.Path)
    b.FromMap(guessed, false)
}

// GuessFromTitle tries to guess some of Book's fields based on Book's Title or
// SubTitle.
func (b *Book) GuessFromTitle() {
	Debug.Printf("Guess subtitle from book's Title '%s'", b.Title)
    guessed := titleGuesser.GuessFrom(b.Title)
    b.FromMap(guessed, false)

	Debug.Printf("Guess information from book's Title '%s'", b.Title)
    guessed = seriesGuesser.GuessFrom(b.Title)
    b.FromMap(guessed, false)

	Debug.Printf("Guess information from book's Sub-Title '%s'", b.Title)
    guessed = seriesGuesser.GuessFrom(b.SubTitle)
    b.FromMap(guessed, false)
}

