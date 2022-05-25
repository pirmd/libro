package main

import (
	"encoding/json"
	"os"

	"github.com/pirmd/libro/book"
	"github.com/pirmd/libro/util"
)

func editBook(editor string, b *book.Book) (*book.Book, error) {
	files, err := book2file(b)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, f := range files {
			_ = os.Remove(f)
		}
	}()

	if err := util.ExecInTTY(editor, files...); err != nil {
		return nil, err
	}

	if len(files) == 0 {
		panic("edit book fail: no temp file containing book's JSON exists")
	}

	if empty, err := util.IsEmptyFile(files[0]); err != nil {
		return nil, err
	} else if empty {
		return nil, nil
	}

	edbook, err := file2book(files[0])
	if err != nil {
		return nil, err
	}

	return edbook, nil
}

func book2file(b *book.Book) ([]string, error) {
	w, err := os.CreateTemp("", "*.json")
	if err != nil {
		return nil, err
	}
	defer func() { _ = w.Close() }()

	prettyJSON := json.NewEncoder(w)
	prettyJSON.SetIndent("", "  ")
	if err := prettyJSON.Encode(struct {
		*book.Book
		Path         string       `json:",omitempty"`
		SimilarBooks []*book.Book `json:",omitempty"`
	}{
		Book: b,
	}); err != nil {
		return nil, err
	}

	// TODO: we call Sync() to capture writing to file errors. An alternative
	// could be to call Close() but it will be redundant with defer (which
	// seems not an issue). Something better could maybe achieved.
	if err := w.Sync(); err != nil {
		return nil, err
	}

	files := []string{w.Name()}

	// During JSON Unmarshal operation, Book.book is not initialized using
	// Book.New() and empty Book.Report will lead to a panic here.
	if b.Report == nil {
		return files, nil
	}

	for _, sb := range b.SimilarBooks {

		filenames, err := book2file(sb)
		if err != nil {
			return nil, err
		}

		files = append(files, filenames...)
	}

	return files, nil
}

func file2book(filename string) (*book.Book, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()

	b := book.New()
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		return nil, err
	}

	return b, nil
}
