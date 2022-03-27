package book

// IsComplete gives an evaluation whether a Book's information seems complete
// enough.
// Objective is to provide a way for the tool to decide whether a human
// intervention is a good idea.
func (b *Book) IsComplete() bool {
	if b.Title == "" {
		Verbose.Printf("warn: Book's Title is empty")
		return false
	}

	if len(b.Authors) == 0 {
		Verbose.Printf("warn: Book has no Author(s)")
		return false
	}

	if b.ISBN == "" {
		Verbose.Printf("warn: Book has no ISBN")
		return false
	}

	if b.Description == "" {
		Verbose.Printf("warn: Book has no Description")
	}

	if b.Publisher == "" {
		Verbose.Printf("warn: Book has no Publisher")
	}

	if b.PublishedDate == "" {
		Verbose.Printf("warn: Book has no PublishedDate")
	}

	if (b.Series != "" && (b.SeriesIndex == 0 || b.SeriesTitle == "")) ||
		(b.SeriesIndex != 0 && (b.Series == "" || b.SeriesTitle == "")) ||
		(b.SeriesTitle != "" && (b.SeriesIndex == 0 || b.Series == "")) {
		Verbose.Printf("warn: Book has incomplete Series information")
		return false
	}

	return true
}
