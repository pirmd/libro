package book

import (
	"strings"

	"github.com/pirmd/libro/book/googlebooks"
)

// FromGooglebooks gets Book's metadata from Googlebooks.
// If search successfully returns a Book with same ISBN, Book's metadata are
// superseded, otherwise the first MaxResult matches are memorized for further
// end-user review only.
func (b *Book) FromGooglebooks(MaxResults int) error {
	google := googlebooks.API{MaxResults: MaxResults}
	found, err := google.SearchVolume(b.toVolumeInfo())
	if err != nil {
		return err
	}

	if len(found) == 0 {
		b.ReportIssue("no match found on Googlebooks")
		return nil
	}

	// TODO: assumes bestMatch should be the one with corresponding ISBN? is it
	// always the case?
	bestMatch := newFromVolumeInfo(found[0])

	if strings.EqualFold(b.ISBN, bestMatch.ISBN) {
		Verbose.Printf("found same book (ISBN: %s) on Googlebooks", b.ISBN)
		Debug.Printf("replace book's metadata with Googlebooks %#v", bestMatch)
		b.ReplaceFrom(bestMatch)
		return nil
	}

	if b.ISBN != "" {
		b.ReportIssue("found book on Googlebooks with different ISBN")
	}

	for _, vi := range found {
		b.ReportSimilarBook(newFromVolumeInfo(vi))
	}
	return nil
}

// toVolumeInfo converts a Book's information into a googlebooks.VolumeInfo.
func (b Book) toVolumeInfo() *googlebooks.VolumeInfo {
	return &googlebooks.VolumeInfo{
		Title:         b.Title,
		SubTitle:      b.SubTitle,
		Language:      b.Language,
		Identifier:    []googlebooks.Identifier{{Type: "ISBN", Identifier: b.ISBN}},
		Authors:       append([]string{}, b.Authors...),
		Subject:       append([]string{}, b.Subject...),
		Description:   b.Description,
		Publisher:     b.Publisher,
		PublishedDate: b.PublishedDate,
		PageCount:     b.PageCount,
	}
}

// newFromVolumeInfo populates Book's information from a googlebooks.VolumeInfo.
func newFromVolumeInfo(vi *googlebooks.VolumeInfo) *Book {
	b := &Book{
		Title:       vi.Title,
		SubTitle:    vi.SubTitle,
		Publisher:   vi.Publisher,
		Description: vi.Description,
		Language:    vi.Language,
		PageCount:   vi.PageCount,
		Subject:     append([]string{}, vi.Subject...),
	}

	b.SetAuthors(vi.Authors)

	isbn := getVolumeInfoISBN(vi)
	b.SetISBN(isbn)

	b.SetPublishedDate(vi.PublishedDate)

	return b
}

func getVolumeInfoISBN(vi *googlebooks.VolumeInfo) (isbn string) {
	for _, id := range vi.Identifier {
		if id.Type == "ISBN" || id.Type == "ISBN_10" || id.Type == "ISBN_13" {
			isbn = id.Identifier
		}

		if len(isbn) == 13 {
			break
		}
	}

	return
}
