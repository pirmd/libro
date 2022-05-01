package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pirmd/libro/book"
	"github.com/pirmd/libro/util"
)

func editBook(editor string, b *book.Book) (*book.Book, error) {
	tmpfile, err := os.CreateTemp("", "*.json")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()

	prettyJSON := json.NewEncoder(tmpfile)
	prettyJSON.SetIndent("", "  ")
	if err := prettyJSON.Encode(b); err != nil {
		return nil, err
	}

	if err := util.ExecInTTY(editor, tmpfile.Name()); err != nil {
		return nil, err
	}

	if _, err := tmpfile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	edbook := book.New()
	if err := json.NewDecoder(tmpfile).Decode(&edbook); err != nil {
		return nil, err
	}

	return edbook, nil
}
