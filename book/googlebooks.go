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
		Debug.Printf("found Googlebook match: %#v", vi)
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
		Debug.Print("no match found on Googlebooks")
		return nil
	}

	bestMatch := newFromVolumeInfo(found[0])

	if strings.EqualFold(b.ISBN, bestMatch.ISBN) {
		Debug.Printf("found same book (ISBN: %s) on Googlebooks", b.ISBN)
		Debug.Printf("replace book's metadata with Googlebooks' one")
		b.ReplaceFrom(bestMatch)
		return nil
	}

	Debug.Printf("found %d similar books on Googlebooks", len(found))
	Debug.Printf("complete book's metadata from Googlebooks' best match")
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
	return &Book{
		Title:         vi.Title,
		Authors:       append([]string{}, vi.Authors...),
		ISBN:          getVolumeInfoISBN(vi),
		SubTitle:      vi.SubTitle,
		Publisher:     vi.Publisher,
		PublishedDate: vi.PublishedDate,
		Description:   vi.Description,
		Language:      vi.Language,
		PageCount:     vi.PageCount,
		Subject:       append([]string{}, vi.Subject...),
	}
}

func getVolumeInfoISBN(vi *googlebooks.VolumeInfo) string {
	for _, id := range vi.Identifier {
		if id.Type == "ISBN" || id.Type == "ISBN_10" || id.Type == "ISBN_13" {
			return id.Identifier
		}
	}

	return ""
}
