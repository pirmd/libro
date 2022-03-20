package book

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

	// reAuthors is a regexp that splits a list of authors.
	reAuthors = regexp.MustCompile(`\s?&\s?`)

	// reCategories is a regexp that splits a list of categories.
	reCategories = regexp.MustCompile(`\s?[&,]\s?`)
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
			return found
		}
	}

	return nil
}

// GuessFromPath tries to guess some of Book's fields based on Book's path.
func (b *Book) GuessFromPath() {
	Debug.Printf("Guess information from book's path '%s'", b.Path)
	b.fromGuesser(pathGuesser.GuessFrom(b.Path))
}

// GuessFromTitle tries to guess some of Book's fields based on Book's Title or
// SubTitle.
func (b *Book) GuessFromTitle() {
	Debug.Printf("Guess subtitle from book's Title '%s'", b.Title)
	b.fromGuesser(titleGuesser.GuessFrom(b.Title))

	Debug.Printf("Guess information from book's Title '%s'", b.Title)
	b.fromGuesser(seriesGuesser.GuessFrom(b.Title))

	Debug.Printf("Guess information from book's Sub-Title '%s'", b.Title)
	b.fromGuesser(seriesGuesser.GuessFrom(b.SubTitle))
}

// fromGuesser update a Book's information according to a (educated?) guesser's
// answer. As we have moderated confidence in the guesser compared to getting
// data from file's Metadata or on-line databases, guesser's result is only
// kept if no previously data is existing. In this case, we only log the
// situation for possibly later improvement of guessers heuristic.
func (b *Book) fromGuesser(guessed map[string]string) {
	Debug.Printf("guessed information: '%+v'", guessed)

	for field, value := range guessed {
		switch f := strings.Title(field); f {
		case "Title":
			if b.Title == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.Title = value
			} else if strings.ToLower(b.Title) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.Title)
			}

		case "SubTitle":
			if b.SubTitle == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.SubTitle = value
			} else if strings.ToLower(b.SubTitle) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.SubTitle)
			}

		case "ShortTitle":
			if b.ShortTitle == "" {
				Verbose.Printf("guess new value: %s => %s.", f, value)
				b.ShortTitle = value
			} else if strings.ToLower(b.ShortTitle) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.ShortTitle)
			}

		case "Authors":
			if auth := reAuthors.Split(value, -1); len(b.Authors) == 0 {
				Verbose.Printf("guess new value: %s = %+v.", f, auth)
				b.Authors = auth
			} else if fmt.Sprint(b.Authors) != fmt.Sprint(auth) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.Authors)
			}

		case "Publisher":
			if b.Publisher == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.Publisher = value
			} else if strings.ToLower(b.Publisher) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.Publisher)
			}

		case "PublishedDate":
			if b.PublishedDate == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.PublishedDate = value
			} else if strings.ToLower(b.PublishedDate) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.PublishedDate)
			}

		case "Description":
			if b.Description == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.Description = value
			} else if strings.ToLower(b.Description) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.Description)
			}

		case "Series":
			if b.Series == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.Series = value
			} else if strings.ToLower(b.Series) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.Series)
			}

		case "SeriesIndex":
			v, err := strconv.ParseFloat(value, 32)
			if err != nil {
				Verbose.Printf("warn: guessed value for '%s' (%s) is not a series number: %v. Ignoring it.", f, value, err)
				break
			}
			if b.SeriesIndex == 0 {
				Verbose.Printf("guess new value: %s = %.1f.", f, v)
				b.SeriesIndex = v
			} else if b.SeriesIndex != v {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.SeriesIndex)
			}

		case "ISBN":
			if b.ISBN == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.ISBN = value
			} else if strings.ToLower(b.ISBN) != strings.ToLower(value) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.ISBN)
			}

		case "Language":
			if b.Language == "" {
				Verbose.Printf("guess new value: %s = %s.", f, value)
				b.Language = value
			} else if strings.ToLower(b.Language) != strings.ToLower(value) {
				Verbose.Printf("war: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.Language)
			}

		case "PageCount":
			v, err := strconv.Atoi(value)
			if err != nil {
				Verbose.Printf("warn: guessed value for '%s' (%s) is not a page number: %v. Ignoring it.", f, value, err)
				break
			}

			if b.PageCount == 0 {
				Verbose.Printf("guess new value: %s = %d.", f, v)
				b.PageCount = v
			} else if b.PageCount != v {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.PageCount)
			}

		case "Categories":
			if cat := reCategories.Split(value, -1); len(b.Categories) == 0 {
				Verbose.Printf("guess new value: %s => %v.", f, cat)
				b.Categories = cat
			} else if fmt.Sprint(b.Categories) != fmt.Sprint(cat) {
				Verbose.Printf("warn: guessed value for '%s' (%s) is different from the existing one (%v). Ignoring it.", f, value, b.Categories)
			}

		default:
			panic("guessed field '" + field + "' is unknown")
		}
	}
}
