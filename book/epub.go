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
	if b.Title == "" {
		Verbose.Printf("warn: no 'Title' found in epub's metadata")
	}

	b.Authors = make([]string, len(mdata.Creator))
	for i, a := range mdata.Creator {
		b.Authors[i] = a.FullName
	}
	if len(b.Authors) == 0 {
		Verbose.Printf("warn: no 'Author' found in epub's metadata")
	}

	if len(mdata.Description) > 0 {
		b.Description = mdata.Description[0]
	}

	b.Subject = append([]string{}, mdata.Subject...)

idloop:
	for _, id := range mdata.Identifier {
		switch {
		case strings.HasPrefix(id.Scheme, "isbn") || strings.HasPrefix(id.Scheme, "ISBN"):
			b.ISBN = id.Value
			break idloop
		case strings.HasPrefix(id.Value, "isbn:"):
			b.ISBN = strings.TrimPrefix(id.Value, "isbn:")
			break idloop
		case strings.HasPrefix(id.Value, "ISBN:"):
			b.ISBN = strings.TrimPrefix(id.Value, "ISBN:")
			break idloop
		}
	}
	if b.ISBN == "" {
		Verbose.Printf("warn: no 'ISBN' found in epub's metadata (%+v)", mdata.Identifier)
	}

	if len(mdata.Publisher) > 0 {
		b.Publisher = mdata.Publisher[0]
	}

	for _, d := range mdata.Date {
		if d.Event == "publication" || d.Event == "" {
			b.PublishedDate = NormalizeDate(d.Stamp)
			break
		}
	}
	if b.PublishedDate == "" {
		Verbose.Printf("warn: no 'publication date' found in epub's metadata (%+v)", mdata.Date)
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
