package book

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hbollon/go-edlib"
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

// String outputs a human understandable description of a SimilarityLevel.
func (lvl SimilarityLevel) String() string {
	return [...]string{"not comparable", "not the same", "maybe the same", "almost the same", "the same"}[lvl]
}

// CompareWith assesses the similarity level between two books with a short
// explanation of the rational.
func (b Book) CompareWith(b1 *Book) (SimilarityLevel, string) {
	isbnLvl := b.compareIdentifierWith(b1)
	nameLvl, nameRational := b.compareNameWith(b1)

	switch isbnLvl {
	case AreTheSame:
		switch nameLvl {
		case AreTheSame:
			return AreTheSame, fmt.Sprintf("ISBN are %s, %s", isbnLvl, nameRational)
		case AreAlmostTheSame, AreNotComparable:
			return AreAlmostTheSame, fmt.Sprintf("ISBN are %s, %s", isbnLvl, nameRational)
		}
		return AreMaybeTheSame, fmt.Sprintf("ISBN are %s, %s", isbnLvl, nameRational)

	case AreNotTheSame:
		switch nameLvl {
		case AreTheSame, AreAlmostTheSame:
			return AreMaybeTheSame, fmt.Sprintf("ISBN are %s, %s", isbnLvl, nameRational)
		default:
			return AreNotTheSame, fmt.Sprintf("ISBN are %s, %s", isbnLvl, nameRational)
		}

	default:
		return nameLvl, fmt.Sprintf("ISBN are %s, %s", isbnLvl, nameRational)
	}
}

func (b Book) compareIdentifierWith(b1 *Book) SimilarityLevel {
	return compareNormalizedISBN(b.ISBN, b1.ISBN)
}

func (b Book) compareNameWith(b1 *Book) (SimilarityLevel, string) {
	lvl := b.compareTitlesWith(b1)

	if lvl >= AreAlmostTheSame {
		if authLvl := b.compareAuthorsWith(b1); authLvl >= AreAlmostTheSame {
			if pubLvl := b.comparePublicationWith(b1); pubLvl >= AreAlmostTheSame {
				return AreTheSame, fmt.Sprintf("Titles are %s, Authors are %s, Publication are %s", lvl, authLvl, pubLvl)
			}
			return AreAlmostTheSame, fmt.Sprintf("Titles are %s, Authors are %s", lvl, authLvl)
		}
		return AreMaybeTheSame, fmt.Sprintf("Titles are %s", lvl)
	}

	return lvl, fmt.Sprintf("Titles are %s", lvl)
}

func (b Book) compareTitlesWith(b1 *Book) SimilarityLevel {
	t, t1 := b.Title, b1.Title

	if b.SubTitle != "" {
		t += " " + b.SubTitle
	}
	if b1.SubTitle != "" {
		t1 += " " + b1.SubTitle
	}

	if t == "" {
		t = b.SeriesTitle
	}
	if t1 == "" {
		t1 = b1.SeriesTitle
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
