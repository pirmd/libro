package book

import (
	"github.com/pirmd/libro/book/googlebooks"
)

// SearchOnGooglebooks search Googlebooks for the Book.
// At most MaxResults Books are returned.
func (b *Book) SearchOnGooglebooks(MaxResults int) ([]*Book, error) {
	api := googlebooks.API{MaxResults: MaxResults}
	found, err := api.SearchVolume(b.toVolumeInfo())
	if err != nil {
		return nil, err
	}

	books := make([]*Book, len(found))
	for i, vi := range found {
		books[i] = newFromVolumeInfo(vi)
	}

	return books, nil
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
	b := New()
	b.Title = vi.Title
	b.SetAuthors(vi.Authors)
	b.SetISBN(getVolumeInfoISBN(vi))
	b.SubTitle = vi.SubTitle
	b.Publisher = vi.Publisher
	b.SetPublishedDate(vi.PublishedDate)
	b.Description = vi.Description
	b.Language = vi.Language
	b.PageCount = vi.PageCount
	b.Subject = append([]string{}, vi.Subject...)

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
