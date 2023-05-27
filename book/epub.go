package book

import (
	"strconv"
	"strings"

	"github.com/pirmd/epub"
)

// NewFromEpub create a Book by populating information out of an EPUB file's
// Metadata.
func NewFromEpub(path string) (*Book, error) {
	b := New()
	b.Path = path

	mdata, err := epub.GetMetadataFromFile(b.Path)
	if err != nil {
		return nil, err
	}

	if len(mdata.Title) > 0 {
		b.Title = mdata.Title[0]
	}

	if len(mdata.SubTitle) > 0 {
		b.SubTitle = mdata.SubTitle[0]
	}

	authors := make([]string, len(mdata.Creator))
	for i, a := range mdata.Creator {
		authors[i] = a.FullName
	}
	b.SetAuthors(authors)

	if len(mdata.Description) > 0 {
		b.SetDescription(mdata.Description[0])
	}

	b.Subject = append([]string{}, mdata.Subject...)

	isbn := getEpubISBN(mdata)
	b.SetISBN(isbn)

	if len(mdata.Publisher) > 0 {
		b.Publisher = mdata.Publisher[0]
	}

	for _, d := range mdata.Date {
		if d.Event == "publication" || d.Event == "" {
			b.SetPublishedDate(d.Stamp)
			break
		}
	}
	if len(mdata.Date) > 0 && b.PublishedDate == "" {
		Debug.Printf("no 'publication date' found in epub's metadata (%+v)", mdata.Date)
	}

	if mdata.Series != "" {
		b.Series = mdata.Series
	}

	if mdata.SeriesIndex != "" {
		v, err := strconv.ParseFloat(mdata.SeriesIndex, 32)
		if err != nil {
			return nil, err
		}
		b.SeriesIndex = v
	}

	if len(mdata.Language) > 0 {
		b.SetLanguage(mdata.Language[0])
	}

	// We extract remaining unused metadata for later improving libro tools.
	for _, meta := range mdata.Meta {
		if meta.Name != "" && meta.Content != "" {
			if meta.Name != "calibre:series" && meta.Name != "calibre:series_index" {
				Debug.Printf("found 'Meta' unused information: %+v", meta)
			}
		}
	}

	return b, nil
}

func getEpubISBN(mdata *epub.Information) (isbn string) {
	for _, id := range mdata.Identifier {
		switch {
		case strings.HasPrefix(id.Scheme, "isbn") || strings.HasPrefix(id.Scheme, "ISBN"):
			isbn = id.Value

		case strings.HasPrefix(id.Value, "urn:isbn:"):
			isbn = strings.TrimPrefix(id.Value, "urn:isbn:")

		case strings.HasPrefix(id.Value, "isbn:"):
			isbn = strings.TrimPrefix(id.Value, "isbn:")

		case strings.HasPrefix(id.Value, "ISBN:"):
			isbn = strings.TrimPrefix(id.Value, "ISBN:")
		}

		// we prefer ISBN_13 over ISBN_10 so if we have it, we're done.
		if len(isbn) == 13 {
			break
		}
	}

	return
}
