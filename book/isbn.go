package book

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// NormalizeISBN returns a cleaned ISBN_13 identifier.
// If isbn is in ISBN_10 format it will be converted to ISBN_13.
// If isbn is in a non-supported format, the 'cleaned' ISBN is returned and an
// error will be raised.
func NormalizeISBN(isbn string) (string, error) {
	if isbn == "" {
		return isbn, nil
	}

	clean := cleanISBN(isbn)

	if len(clean) == 13 {
		return clean, nil
	}

	isbn13, err := toISBN13(clean)
	if err != nil {
		return clean, err
	}

	return isbn13, nil
}

// cleanISBN returns an ISBN identifier without any separator or blank.
func cleanISBN(isbn string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == 'X' {
			return r
		}
		return -1
	}, isbn)
}

// toISBN13 tries to convert an ISBN_10 to an ISBN_13.
// isbn10 is supposed to be a "cleaned" ISBN (only digits, no '-' or things
// like that)
func toISBN13(isbn10 string) (string, error) {
	if len(isbn10) != 10 {
		return "", fmt.Errorf("convert to ISBN_13 failed: %s is not a suitable ISBN_10", isbn10)
	}

	isbn13 := "978" + isbn10[:9]

	checkdigit, err := calcISBN13checkdigit(isbn13)
	if err != nil {
		return "", err
	}

	return isbn13 + checkdigit, nil
}

// calcISBN13checkdigit calculates last check-digit of an ISBN_13.
// isbn13 should be a 'cleaned' isbn13.
// if isbn13 is provided with a 13th digit, check-digit calculation outcome
// will be compared to this 13th digit and an error will be raised if they are
// different.
func calcISBN13checkdigit(isbn13 string) (string, error) {
	if len(isbn13) != 12 && len(isbn13) != 13 {
		return "", fmt.Errorf("ISBN_13 check-digit calculation failed: %s is not a suitable ISBN_13", isbn13)
	}

	// from: https://en.wikipedia.org/wiki/International_Standard_Book_Number#ISBN-13_check_digit_calculation
	var sum int
	for i, c := range isbn13[:12] {
		sum += int(c-'0') * (1 + 2*(i%2))
	}

	digit := 10 - sum%10

	if digit == 10 {
		digit = 0
	}
	checkdigit := strconv.Itoa(digit)

	if len(isbn13) == 13 && isbn13[12:] != checkdigit {
		return "", fmt.Errorf("mismatched check-digit between actual ISBN_13 (%s) and calculated one (%s)", isbn13[12:], checkdigit)
	}
	return checkdigit, nil
}
