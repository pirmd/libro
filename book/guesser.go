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
	// Reminder for guessers in this section: Order is important as only first
	// match is considered, so it needs to be defined from the more specific to
	// the more general capture logic.

	// pathGuessers is a collection of regexp to extract information from a
	// Book's filename.
	pathGuessers = []*regexp.Regexp{
		// parent/folder/<Authors> - [<Series> <SeriesIndex>] - <SeriesTitle> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s-\s\[(?P<Series>.+)\s(?P<SeriesIndex>\d+)\]\s-\s(?P<SeriesTitle>.+?)\s\[(?P<Language>.+)\]\.(?:.+)$`),
		// parent/folder/<Authors> - <Title> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s-\s(?P<Title>.+?)\s\[(?P<Language>.+)\]\.(?:.+)$`),
	}

	// seriesGuessers is a collection of regexp to extract series information
	// from a Book's title or subtitle.
	seriesGuessers = []*regexp.Regexp{
		// <SeriesTitle> (<Series> n°<SeriesIndex>)
		regexp.MustCompile(`^(?P<SeriesTitle>.+)\s\p{Ps}(?P<Series>.+?)\s` + reSeriesIndex + `\p{Pe}$`),
		// <SeriesTitle> - <Series> n°<SeriesIndex>
		regexp.MustCompile(`^(?P<SeriesTitle>.+?)\s\p{Pd}\s(?P<Series>.+?)\s` + reSeriesIndex + `$`),
		// <Series> n°<SeriesIndex>
		regexp.MustCompile(`^(?P<Series>.+?)\s` + reSeriesIndex + `$`),
		// Book <SeriesIndex> of <Series>
		regexp.MustCompile(`^Book\s(?P<SeriesIndex>\d+)\sof\s(?P<Series>.+)$`),
	}

	// titleCleaners is a collection of regexp that pre-processes bad-formatted Titles
	titleCleaners = []*regexp.Regexp{
		// <Title> / <SubTitle>
		regexp.MustCompile(`^(?P<Title>.+)\s/\s(?P<SubTitle>.+)$`),
		// <Title> (French Edition)
		regexp.MustCompile(`^(?P<Title>.+)\s\p{Ps}(?i:.+\sedition)\p{Pe}$`),
	}
)

// GuessFromFilename tries to guess Book's information based on its filename.
func (b *Book) GuessFromFilename() error {
	Debug.Printf("Guess information from book's path '%s'", b.Path)
	if err := b.guess(b.Path, pathGuessers...); err != nil {
		return err
	}

	return nil
}

// GuessFromMetadata tries to guess Book's information based on known
// attributes (like Book's Title).
func (b *Book) GuessFromMetadata() error {
	Debug.Printf("Guess series information from book's Title '%s'", b.Title)
	if err := b.guess(b.Title, seriesGuessers...); err != nil {
		return err
	}

	if b.SubTitle != "" {
		Debug.Printf("Guess series information from book's Sub-Title '%s'", b.SubTitle)
		if err := b.guess(b.SubTitle, seriesGuessers...); err != nil {
			return err
		}
	}

	return nil
}

// CleanMetadata cleans Book's metadata.
func (b *Book) CleanMetadata() error {
	Debug.Printf("Clean book's Title '%s'", b.Title)
	if err := b.clean(b.Title, titleCleaners...); err != nil {
		return err
	}

	return nil
}

// guess extracts new Book's attributes from a string by applying a list of
// Regexp.
// Regexp guesses new attribute's value using capturing group whose name shall
// correspond to the attribute to create. Unknown attribute name will raise an
// error.
// Regexp are run in guessers order and only first match is returned.
func (b *Book) guess(s string, guessers ...*regexp.Regexp) error {
	for _, re := range guessers {
		guessed := submatchAsMap(s, re)
		if guessed != nil {
			Debug.Printf("guessed information: '%+v'", guessed)
			return b.CompleteFromMap(guessed)
		}
	}

	Debug.Printf("no match found")
	return nil
}

// clean rewrites Book's attributes by applying a list of Regexp.
// Regexp guesses new attribute's value using capturing group whoses nama
// correspond to the attribute to update or to create. Unknown attribute name
// will raise an error.
// Regexp are run in the cleaners order.
func (b *Book) clean(s string, cleaners ...*regexp.Regexp) error {
	for _, re := range cleaners {
		cleaned := submatchAsMap(s, re)
		if cleaned != nil {
			Debug.Printf("clean information: '%+v'", cleaned)
			if err := b.ReplaceFromMap(cleaned); err != nil {
				return err
			}
		}
	}

	Debug.Printf("no match found")
	return nil
}
