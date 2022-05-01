package book

import (
	"strings"

	"github.com/pirmd/libro/book/googlebooks"
)

// SearchGooglebooks search Googlebooks for same or similar Books.
func (b Book) SearchGooglebooks() ([]*Book, error) {
	found, err := googlebooks.SearchVolume(b.toVolumeInfo())
	if err != nil {
		return nil, err
	}

	var res []*Book
	for _, vi := range found {
		Debug.Printf("found Googlebooks match: %#v", vi)
		res = append(res, newFromVolumeInfo(vi))
	}

	return res, nil
}

// FromGooglebooks gets Book's metadata from Googlebooks.
// If searching Googlebooks successfully found a Book with same ISBN, Book's
// metadata are superseded by Googlebooks' one, otherwise the best match is use
// to complete Book's metada.
func (b *Book) FromGooglebooks() error {
	found, err := googlebooks.SearchVolume(b.toVolumeInfo())
	if err != nil {
		return err
	}

	if len(found) == 0 {
		Verbose.Print("warn: no match found on Googlebooks")
		return nil
	}

	bestMatch := newFromVolumeInfo(found[0])

	if strings.EqualFold(b.ISBN, bestMatch.ISBN) {
		Verbose.Printf("found same book (ISBN: %s) on Googlebooks", b.ISBN)
		Debug.Printf("replace book's metadata with Googlebooks %#v", bestMatch)
		b.ReplaceFrom(bestMatch)
		return nil
	}

	if b.ISBN != "" {
		Verbose.Printf("warn: found book on Googlebooks with different ISBN (%s != %s), ignore it", b.ISBN, bestMatch.ISBN)
		return nil
	}

	Verbose.Printf("found %d similar books on Googlebooks", len(found))
	Debug.Printf("complete book's metadata from Googlebooks %#v", bestMatch)
	b.CompleteFrom(bestMatch)
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
