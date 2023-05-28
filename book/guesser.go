package book

import (
	"io"
	"io/fs"
	"regexp"

	"github.com/pirmd/epub"

	"github.com/pirmd/libro/book/htmlutil"
)

const (
	// reSeriesIndex is a regexp that captures a series index.
	reSeriesIndex = `(?i:Series |Volume |Vol.* |Part |Livre |Tome |n°|#|T|)(?P<SeriesIndex>\d{1,3})`

	// reISBN is a regexp aiming at capturing ISBN-like indication in text. It does not
	// aim at validating an ISBN, it can even return too short or too long results.
	// Libro should rely on NormalizeISBN step to make sure detected ISBN is valid.
	reISBN = `(?P<ISBN>(?:97[89][\d\p{Zs}\p{Pd}]{10,14})|(?:[\d][\d\p{Zs}\p{Pd}]{8,11}[\dxX]))`

	// reExt is a regexp aiming at capturing any 'reasonable' filename extension.
	reExt = `\.[\w]+$`

	// reLang is a regexp aiming at capturing any 'reasonable' language identifiers.
	reLang = `\s*\p{Ps}(?P<Language>[_-a-zA-Z]{2,5})\p{Pe}`
)

var (
	// reSeriesWithIndex is a regexp that captures a Series and its index
	reSeriesWithIndex = `(?P<Series>.+?)(?:\s?\p{Pd}\s?|,\s|\s)` + reSeriesIndex

	// Reminder for guessers in this section: order is important as only first
	// match is considered, so it needs to be defined from the more specific to
	// the more general capture logic.

	// pathGuessers is a collection of regexp to extract information from a
	// Book's filename.
	pathGuessers = []*regexp.Regexp{
		// parent/folder/<Authors> - [<Series> <SeriesIndex>] - <SeriesTitle> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s\p{Pd}\s\p{Ps}` + reSeriesWithIndex + `\p{Pe}\s\p{Pd}\s(?P<SeriesTitle>.+?)` + reLang + reExt),
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s\p{Pd}\s\p{Ps}` + reSeriesWithIndex + `\p{Pe}\s\p{Pd}\s(?P<SeriesTitle>.+?)` + reExt),
		// parent/folder/<Authors> - <Series> <SeriesIndex> - <SeriesTitle> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s\p{Pd}\s` + reSeriesWithIndex + `\s\p{Pd}\s(?P<SeriesTitle>.+?)` + reLang + reExt),
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s\p{Pd}\s` + reSeriesWithIndex + `\s\p{Pd}\s(?P<SeriesTitle>.+?)` + reExt),
		// parent/folder/<Authors> - <Title> [<Language>].epub
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s\p{Pd}\s(?P<Title>.+?)` + reLang + reExt),
		regexp.MustCompile(`^(?:.*/)?(?P<Authors>.+)\s\p{Pd}\s(?P<Title>.+?)` + reExt),
	}

	// seriesGuessers is a collection of regexp to extract series information
	// from a Book's title or subtitle.
	seriesGuessers = []*regexp.Regexp{
		// Book <SeriesIndex> of <Series>
		regexp.MustCompile(`^Book\s(?P<SeriesIndex>\d+)\sof\s(?P<Series>.+)$`),
		// <SeriesTitle> (<Series> <SeriesIndex>)
		regexp.MustCompile(`^(?P<SeriesTitle>.+)\s\p{Ps}` + reSeriesWithIndex + `\p{Pe}$`),
		// <SeriesTitle> - <Series> <SeriesIndex>
		regexp.MustCompile(`^(?P<SeriesTitle>.+?)\s\p{Pd}\s` + reSeriesWithIndex + `$`),
		// [<Series> <SeriesIndex>] <SeriesTitle>
		regexp.MustCompile(`^\p{Ps}` + reSeriesWithIndex + `\p{Pe}\s(?P<SeriesTitle>.+)$`),
		// <Series> <SeriesIndex> - <SeriesTitle>
		regexp.MustCompile(`^` + reSeriesWithIndex + `(?:\s\p{Pd}|,|\s?:)\s(?P<SeriesTitle>.+)$`),
		// <Series> (<SeriesIndex>) - <SeriesTitle>
		regexp.MustCompile(`^(?P<Series>.+?)\s\p{Ps}` + reSeriesIndex + `\p{Pe}(?:\s\p{Pd}|,)*\s(?P<SeriesTitle>.+)$`),
		// <Series> <SeriesIndex>
		regexp.MustCompile(`^` + reSeriesWithIndex + `$`),
		// <SeriesIndex> - <SeriesTitle>
		regexp.MustCompile(`^` + reSeriesIndex + `\s*[.,-]\s*(?P<SeriesTitle>.+)$`),
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
		regexp.MustCompile(`^(?P<Title>.+)\s\p{Ps}(?:.*(?i:[eé]dition|litt[eé]rature).*|\p{Lu}+|[0-9]{4})\p{Pe}$`),
	}
)

// NewFromFilename creates a Book whose information are guessed from its filename.
func NewFromFilename(path string) (*Book, error) {
	return guess(path, pathGuessers...)
}

// NewFromContent creates a Book whose information are guessed from its Content.
func NewFromContent(path string) (*Book, error) {
	return grep(path, contentGuesser)
}

// GuessFromMetadata tries to guess Book's information based on known
// attributes (like Book's Title).
func (b *Book) GuessFromMetadata() error {
	if b.Title != "" {
		Debug.Printf("guess Series from Title '%s'", b.Title)
		if err := b.guess(b.Title, seriesGuessers...); err != nil {
			return err
		}
	}

	if b.SubTitle != "" {
		Debug.Printf("guess Series from Sub-Title '%s'", b.SubTitle)
		if err := b.guess(b.SubTitle, seriesGuessers...); err != nil {
			return err
		}
	}

	return nil
}

// CleanMetadata cleans Book's metadata.
func (b *Book) CleanMetadata() error {
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
	guessedBook, err := guess(s, guessers...)
	if err != nil {
		return err
	}

	if guessedBook != nil {
		b.CompleteFrom(guessedBook)
	}

	return nil
}

// guess extracts new Book's attributes from a string by applying a list of
// Regexp.
// Regexp guesses new attribute's value using capturing group whose name shall
// correspond to the attribute to create. Unknown attribute name will raise an
// error.
// Regexp are run in guessers order and only first match is returned.
func guess(s string, guessers ...*regexp.Regexp) (*Book, error) {
	for _, re := range guessers {
		if guessed := reFindStringSubmatchAsMap(s, re); guessed != nil {
			Debug.Printf("guessed information: '%+v'", guessed)
			return NewFromMap(guessed)
		}
	}

	return nil, nil
}

// grep extracts new Book's attributes from its content by applying a Regexp.
// Regexp guesses new attribute's value using capturing group whose name shall
// correspond to the attribute to create. Unknown attribute name will raise an
// error.
// Several matches for the same attribute can be returned, management of
// inconsistent values is left to Book.CompleteFromMap logic, eventually
// reporting to end-user such situation.
func grep(path string, re *regexp.Regexp) (*Book, error) {
	// TODO: I'm quite 'defensive' here as I capture every matches and report
	// possible inconsistent values. This can maybe be removed later one once
	// better confident in the heuristic so that we can just stop on the first
	// match.
	var found []map[string]string

	if err := epub.WalkReadingContent(path, func(r io.Reader, fi fs.FileInfo) error {
		rawr, err := htmlutil.GetRawTextFromHTML(r)
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
		return nil, err
	}

	if len(found) == 0 {
		Debug.Printf("no match found")
		return nil, nil
	}

	b, err := NewFromMap(found[0])
	if err != nil {
		return nil, err
	}

	for _, f := range found[1:] {
		if err := b.CompleteFromMap(f); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// clean rewrites Book's attributes by applying a list of Regexp.
// Regexp guesses new attribute's value using capturing group whose name
// correspond to the attribute to update or to create. Unknown attribute name
// will raise an error.
// Regexp are run in the cleaners order.
func (b *Book) clean(s string, cleaners ...*regexp.Regexp) error {
	for _, re := range cleaners {
		if cleaned := reFindStringSubmatchAsMap(s, re); cleaned != nil {
			Debug.Printf("clean information: '%+v'", cleaned)
			if err := b.ReplaceFromMap(cleaned); err != nil {
				return err
			}
		}
	}

	return nil
}
