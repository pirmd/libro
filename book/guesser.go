package book

import (
	"io"
	"io/fs"
	"regexp"

	"github.com/pirmd/epub"
)

const (
	// seriesKey lists usual keywords that introduce a specification of a
	// series index.
	reSeriesIndex = `(?i:Series |Volume |Vol |Vol. |Part |n°|#|)(?P<SeriesIndex>\d+)`

	// reISBN is a regexp aiming at capturing ISBN-like indication in text. It does not
	// aim at validating an ISBN, it can even return too short or too long results.
	// Libro should rely on NormalizeISBN step to make sure detected ISBN is valid.
	reISBN = `(?P<ISBN>(?:97[89][\d\p{Zs}\p{Pd}]{10,14})|(?:[\d][\d\p{Zs}\p{Pd}]{8,11}[\dxX]))`
)

var (
	// Reminder for guessers in this section: order is important as only first
	// match is considered, so it needs to be defined from the more specific to
	// the more general capture logic.

	// pathGuessers is a collection of regexp to extract information from a
	// Book's filename.
	pathGuessers = []*regexp.Regexp{
		// parent/folder/<Authors> - [<Series> <SeriesIndex>] - <SeriesTitle> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s-\s\[(?P<Series>.+)\s(?P<SeriesIndex>\d+)\]\s-\s(?P<SeriesTitle>.+?)\s\[(?P<Language>.+)\]\.(?:.+)$`),
		// parent/folder/<Authors> - <Series> <SeriesIndex> - <SeriesTitle> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s-\s(?P<Series>.+)\s(?P<SeriesIndex>\d+)\s-\s(?P<SeriesTitle>.+?)\s\[(?P<Language>.+)\]\.(?:.+)$`),
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

	// contentGuesser is a regexp that extracts information from a Book's content.
	// ISBN: <isbn> ou EAN: <isbn>
	contentGuesser = regexp.MustCompile(
		`(?:(?:ISBN)|(?:EAN)).*?\p{Zs}?:?\p{Zs}?` + reISBN,
	)

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

// GuessFromContent tries to guess Book's information from the Book's content.
func (b *Book) GuessFromContent() error {
	Debug.Print("Guess ISBN information from book's Content")
	if err := b.grep(contentGuesser); err != nil {
		return err
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
		guessed := reFindStringSubmatchAsMap(s, re)
		if guessed != nil {
			Debug.Printf("guessed information: '%+v'", guessed)
			return b.CompleteFromMap(guessed)
		}
	}

	Debug.Printf("no match found")
	return nil
}

// grep extracts new Book's attributes from a Reader by applying a Regexp.
// Regexp guesses new attribute's value using capturing group whose name shall
// correspond to the attribute to create. Unknown attribute name will raise an
// error.
// Several matches for the same attribute can be returned, management of
// inconsistent values is left to Book.CompleteFromMap logic, eventually
// reporting to end-user such situation.
func (b *Book) grep(re *regexp.Regexp) error {
	// TODO: I'm quite 'defensive' here as I capture every matches and report
	// possible inconsistent values. This can maybe be removed later one once
	// better confident in the heuristic so that we can just stop on the first
	// match.
	var found []map[string]string

	if err := epub.WalkReadingContent(b.Path, func(r io.Reader, fi fs.FileInfo) error {
		rawr, err := getRawTextFromHTML(r)
		if err != nil {
			return err
		}

		matches := reFindReaderSubmatchAsMap(rawr, re)
		if matches != nil {
			Debug.Printf("found information in %s: '%+v'", fi.Name(), matches)
			found = append(found, matches...)
		}

		return nil
	}); err != nil {
		return err
	}

	if len(found) == 0 {
		Debug.Printf("no match found")
	}

	for _, f := range found {
		if err := b.CompleteFromMap(f); err != nil {
			return err
		}
	}
	return nil
}

// clean rewrites Book's attributes by applying a list of Regexp.
// Regexp guesses new attribute's value using capturing group whoses nama
// correspond to the attribute to update or to create. Unknown attribute name
// will raise an error.
// Regexp are run in the cleaners order.
func (b *Book) clean(s string, cleaners ...*regexp.Regexp) error {
	for _, re := range cleaners {
		cleaned := reFindStringSubmatchAsMap(s, re)
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
