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

	authors := make([]string, len(mdata.Creator))
	for i, a := range mdata.Creator {
		authors[i] = a.FullName
	}
	b.SetAuthors(authors)

	if len(mdata.Description) > 0 {
		b.Description = mdata.Description[0]
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

	for _, meta := range mdata.Meta {
		switch meta.Name {
		case "calibre:series":
			b.Series = meta.Content

		case "calibre:series_index":
			v, err := strconv.ParseFloat(meta.Content, 32)
			if err != nil {
				return nil, err
			}
			b.SeriesIndex = v

		default:
			if meta.Name != "" || meta.Content != "" {
				Debug.Printf("found 'Meta' unused information: %+v", meta)
			}
		}
	}

	return b, nil
}

func getEpubISBN(mdata *epub.Metadata) (isbn string) {
	for _, id := range mdata.Identifier {
		switch {
		case strings.HasPrefix(id.Scheme, "isbn") || strings.HasPrefix(id.Scheme, "ISBN"):
			isbn = id.Value

		case strings.HasPrefix(id.Value, "isbn:"):
			isbn = strings.TrimPrefix(id.Value, "isbn:")

		case strings.HasPrefix(id.Value, "ISBN:"):
			isbn = strings.TrimPrefix(id.Value, "ISBN:")
		}

		if len(isbn) == 13 {
			break
		}
	}

	return
}
