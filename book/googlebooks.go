package book

import (
	"fmt"
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

	if strings.EqualFold(b.ISBN, getVolumeInfoISBN(found[0])) {
		Debug.Printf("found same book (ISBN: %s) on Googlebooks", b.ISBN)
		Debug.Printf("replace book's metadata with Googlebooks' one")
		b.fromVolumeInfo(found[0], true)
		return nil
	}

	Debug.Printf("found %d similar books on Googlebooks", len(found))
	Debug.Printf("complete book's metadata from Googlebooks' best match")
	b.fromVolumeInfo(found[0], false)
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

// fromVolumeInfo populates Book's information from a googlebooks.VolumeInfo.
func (b *Book) fromVolumeInfo(vi *googlebooks.VolumeInfo, override bool) {
	if vi.Title != "" {
		if b.Title != "" && strings.EqualFold(b.Title, vi.Title) {
			Debug.Printf("new Book's Title is different from the existing one (%s != %s)", vi.Title, b.Title)
		}
		if override || b.Title == "" {
			Verbose.Printf("sets new Book's value: Title = %s", b.Title)
			b.Title = vi.Title
		}
	}

	if vi.SubTitle != "" {
		if b.SubTitle != "" && strings.EqualFold(b.SubTitle, vi.SubTitle) {
			Debug.Printf("new Book's SubTitle is different from the existing one (%s != %s)", vi.SubTitle, b.SubTitle)
		}
		if override || b.SubTitle == "" {
			Verbose.Printf("sets new Book's value: SubTitle = %s", b.SubTitle)
			b.SubTitle = vi.SubTitle
		}
	}

	if len(vi.Authors) > 0 {
		if len(b.Authors) != 0 && strings.EqualFold(fmt.Sprint(b.Authors), fmt.Sprint(vi.Authors)) {
			Debug.Printf("new Book's Authors is different from the existing one (%v != %v)", vi.Authors, b.Authors)
		}
		if override || len(b.Authors) == 0 {
			Verbose.Printf("sets new Book's value: Authors = %v", vi.Authors)
			b.Authors = append([]string{}, vi.Authors...)
		}
	}

	if vi.Publisher != "" {
		if b.Publisher != "" && strings.EqualFold(b.Publisher, vi.Publisher) {
			Debug.Printf("new Book's Publisher is different from the existing one (%s != %s)", vi.Publisher, b.Publisher)
		}
		if override || b.Publisher == "" {
			Verbose.Printf("sets new Book's value: Publisher = %s", vi.Publisher)
			b.Publisher = vi.Publisher
		}
	}

	if vi.PublishedDate != "" {
		if b.PublishedDate != "" && strings.EqualFold(b.PublishedDate, vi.PublishedDate) {
			Debug.Printf("new Book's PublishedDate is different from the existing one (%s != %s)", vi.PublishedDate, b.PublishedDate)
		}
		if override || b.PublishedDate == "" {
			Verbose.Printf("sets new Book's value: PublishedDate = %s", vi.PublishedDate)
			b.PublishedDate = vi.PublishedDate
		}
	}

	if vi.Description != "" {
		if b.Description != "" && strings.EqualFold(b.Description, vi.Description) {
			Debug.Printf("new Book's Description is different from the existing one (%s != %s)", vi.Description, b.Description)
		}
		if override || b.Description == "" {
			Verbose.Printf("sets new Book's value: Description = %s", vi.Description)
			b.Description = vi.Description
		}
	}

	isbn := getVolumeInfoISBN(vi)
	if isbn != "" {
		if b.ISBN != "" && strings.EqualFold(b.ISBN, isbn) {
			Debug.Printf("new Book's ISBN is different from the existing one (%s != %s)", isbn, b.ISBN)
		}
		if override || b.ISBN == "" {
			Verbose.Printf("sets new Book's value: ISBN = %s", isbn)
			b.ISBN = isbn
		}
	}

	if vi.Language != "" {
		if b.Language != "" && strings.EqualFold(b.Language, vi.Language) {
			Debug.Printf("new Book's Language is different from the existing one (%s != %s)", vi.Language, b.Language)
		}
		if override || b.Language == "" {
			Verbose.Printf("sets new Book's value: Language = %s", vi.Language)
			b.Language = vi.Language
		}
	}

	if vi.PageCount != 0 {
		if b.PageCount != 0 && b.PageCount != vi.PageCount {
			Debug.Printf("new Book's PageCount is different from the existing one (%d != %d)", vi.PageCount, b.PageCount)
		}
		if override || b.PageCount == 0 {
			Verbose.Printf("sets new value: PageCount = %d", vi.PageCount)
			b.PageCount = vi.PageCount
		}
	}

	if len(vi.Subject) > 0 {
		if len(b.Subject) != 0 && strings.EqualFold(fmt.Sprint(b.Subject), fmt.Sprint(vi.Subject)) {
			Debug.Printf("new Book's Subject is different from the existing one (%v != %v)", vi.Subject, b.Subject)
		}
		if override || len(b.Subject) == 0 {
			Verbose.Printf("sets new Book's value: Subject = %v", vi.Subject)
			b.Subject = append([]string{}, vi.Subject...)
		}
	}
}

// newFromVolumeInfo populates Book's information from a googlebooks.VolumeInfo.

func newFromVolumeInfo(vi *googlebooks.VolumeInfo) *Book {
	b := New()
	b.fromVolumeInfo(vi, true)
	return b
}

func getVolumeInfoISBN(vi *googlebooks.VolumeInfo) string {
	for _, id := range vi.Identifier {
		if id.Type == "ISBN" || id.Type == "ISBN_10" || id.Type == "ISBN_13" {
			return id.Identifier
		}
	}

	return ""
}
