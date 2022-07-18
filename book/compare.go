package book

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hbollon/go-edlib"
)

const (
	// isAreTheSameThreshold is the minimum distance above which two strings are
	// AreTheSame.
	isAreTheSameThreshold = 0.9
	// isAreAlmostTheSameThreshold is the minimum distance above which two strings are
	// AreAlmostTheSame.
	isAreAlmostTheSameThreshold = 0.8
	// isAreMaybeTheSameThreshold is the minimum distance above which two strings are
	// AreMaybeTheSame.
	isAreMaybeTheSameThreshold = 0.7
)

// SimilarityLevel indicate the similarity level between two elements.
type SimilarityLevel int

const (
	// AreNotComparable indicates that two elements are not comparable.
	AreNotComparable SimilarityLevel = iota
	// AreNotTheSame indicates that two elements are different.
	AreNotTheSame
	// AreMaybeTheSame indicates that two elements are maybe the same.
	AreMaybeTheSame
	// AreAlmostTheSame indicates that elements are almost the same.
	AreAlmostTheSame
	// AreTheSame indicates that elements are the same.
	AreTheSame
)

// String outputs a human understandable description of a SimilarityLevel.
func (lvl SimilarityLevel) String() string {
	return [...]string{"same", "almost the same", "maybe the same", "not the same", "not comparable"}[lvl]
}

// CompareWith assesses the similarity level between two books with a short
// explanation of the rational.
func (b Book) CompareWith(b1 *Book) (SimilarityLevel, string) {
	//TODO: fine-tunes rational by detailing 'different Name' or 'similar
	//Names' cause, possibly by showing the Name or getting rational from
	//compareNameWith

	Debug.Printf("compare %#v with %#v", b, b1)
	switch b.compareIdentifierWith(b1) {
	case AreTheSame:
		switch b.compareNameWith(b1) {
		case AreTheSame:
			return AreTheSame, "same ISBN and Names"
		case AreAlmostTheSame:
			return AreAlmostTheSame, "same ISBN and similar Names"
		case AreNotComparable:
			return AreAlmostTheSame, "same ISBN but not comparable Names"
		}
		return AreMaybeTheSame, "Same ISBN but different Name"

	case AreNotTheSame:
		switch l := b.compareNameWith(b1); l {
		case AreTheSame, AreAlmostTheSame:
			return AreMaybeTheSame, fmt.Sprintf("different ISBN (%s vs. %s) and Names are %s", b.ISBN, b1.ISBN, l)
		default:
			return AreNotTheSame, fmt.Sprintf("different ISBN (%s vs %s) and Names are %s", b.ISBN, b1.ISBN, l)
		}

	default:
		l := b.compareNameWith(b1)
		return l, fmt.Sprintf("ISBN are not comparable, Names are %s", l)
	}
}

func (b Book) compareIdentifierWith(b1 *Book) SimilarityLevel {
	lvl := compareNormalizedISBN(b.ISBN, b1.ISBN)
	Debug.Printf("book's ISBN are %s", lvl)
	return lvl
}

func (b Book) compareNameWith(b1 *Book) SimilarityLevel {
	lvl := b.compareTitlesWith(b1)

	if lvl >= AreAlmostTheSame {
		if authLvl := b.compareAuthorsWith(b1); authLvl >= AreAlmostTheSame {
			if pubLvl := b.comparePublicationWith(b1); pubLvl >= AreAlmostTheSame {
				Debug.Printf("books' Title are %s, books' Authors are %s, books' Publication are %s -> are the same", lvl, authLvl, pubLvl)
				return AreTheSame
			}
			Debug.Printf("books' Title are %s, books' Authors are %s -> are almost same", lvl, authLvl)
			return AreAlmostTheSame
		}
		Debug.Printf("books' Title are %s -> are maybe the same", lvl)
		return AreMaybeTheSame
	}

	Debug.Printf("books' Title are %s -> are maybe the same", lvl)
	return lvl
}

func (b Book) compareTitlesWith(b1 *Book) SimilarityLevel {
	t, t1 := b.Title, b1.Title

	if b.SubTitle != "" {
		t += " " + b.SubTitle
	}
	if b1.SubTitle != "" {
		t1 += " " + b1.SubTitle
	}

	return compareStrings(t, t1)
}

func (b Book) compareAuthorsWith(b1 *Book) SimilarityLevel {
	return compareLists(b.Authors, b1.Authors)
}

func (b Book) comparePublicationWith(b1 *Book) SimilarityLevel {
	lvl := b.comparePublisherWith(b1)
	if lvl >= AreAlmostTheSame {
		if b.comparePublishedDateWith(b1) >= AreAlmostTheSame {
			return AreTheSame
		}
		return AreAlmostTheSame
	}

	return lvl
}

func (b Book) comparePublisherWith(b1 *Book) SimilarityLevel {
	return compareStrings(b.Publisher, b1.Publisher)
}

func (b Book) comparePublishedDateWith(b1 *Book) SimilarityLevel {
	return compareNormalizedDates(b.PublishedDate, b1.PublishedDate)
}

func (b Book) compareTitleWith(b1 *Book) SimilarityLevel {
	return compareStrings(b.Title, b1.Title)
}

func (b Book) compareSubTitleWith(b1 *Book) SimilarityLevel {
	return compareStrings(b.SubTitle, b1.SubTitle)
}

func (b Book) compareSubjectWith(b1 *Book) SimilarityLevel {
	return compareLists(b.Subject, b1.Subject)
}

// compareStrings compares two strings considering their Jaro-Winkler distance
func compareStrings(s1, s2 string) SimilarityLevel {
	if s1 == "" || s2 == "" {
		return AreNotComparable
	}

	dist, err := edlib.StringsSimilarity(s1, s2, edlib.JaroWinkler)
	if err != nil {
		return AreNotComparable
	}

	switch {
	case dist > isAreTheSameThreshold:
		return AreTheSame
	case dist > isAreAlmostTheSameThreshold:
		return AreAlmostTheSame
	case dist > isAreMaybeTheSameThreshold:
		return AreMaybeTheSame
	}

	return AreNotTheSame
}

// compareLists compare two lists of strings, without considering order.
func compareLists(l1, l2 []string) SimilarityLevel {
	// TODO: improve heuristic to differentiate completely different lists from
	// lists that are close like toto vs tata, toto -> AreAlmostTheSame
	sorted1, sorted2 := append([]string{}, l1...), append([]string{}, l2...)
	sort.Strings(sorted1)
	sort.Strings(sorted2)
	return compareStrings(fmt.Sprint(sorted1), fmt.Sprint(sorted2))
}

// compareNormalizedISBN compares two already 'normalized' ISBN.
func compareNormalizedISBN(isbn1, isbn2 string) SimilarityLevel {
	switch {
	case isbn1 == "" || isbn2 == "":
		return AreNotComparable

	case len(isbn1) != len(isbn2):
		return AreNotComparable

	case isbn1 == isbn2:
		return AreTheSame

	default:
		return AreNotTheSame
	}
}

// compareNormalizedDates compares two already 'normalized' date.
func compareNormalizedDates(date1, date2 string) SimilarityLevel {
	if date1 == "" || date2 == "" {
		return AreNotComparable
	}

	if date1 == date2 {
		return AreTheSame
	}

	d1, d2 := date1, date2
	if len(date1) < len(date2) {
		d1, d2 = date2, date1
	}

	if strings.HasPrefix(d1, d2) {
		return AreAlmostTheSame
	}

	return AreNotTheSame
}
