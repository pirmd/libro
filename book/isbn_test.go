package book

import (
	"testing"
)

func TestCalcISBN13checkdigit(t *testing.T) {
	testCases := []string{
		"9780002712712",
		"9780002712835",
		"9780002713023",
		"9780002713122",
		"9780002713276",
		"9780002713306",
		"9780002713320",
		"9780002714648",
		"9780002714686",
		"9780002715065",
		"9780002715096",
		"9780002715102",
		"9780002712095",
		"9780002712149",
		"9780002712170",
		"9780002712187",
	}

	for _, tc := range testCases {
		want := tc[12:]
		got, err := calcISBN13checkdigit(tc[:12])
		if err != nil {
			t.Errorf("Fail to calculate check-digit for %s: %v", tc, err)
			continue
		}

		if want != got {
			t.Errorf("fail to calculate check-digit for %s. Got: %v. Want: %v", tc, got, want)
		}
	}
}

func TestToISBN13(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{out: "9782211238434", in: "2211238432"},
		{out: "9780486285405", in: "0486285405"},
		{out: "9782017064657", in: "2017064653"},
	}

	for _, tc := range testCases {
		got, err := toISBN13(tc.in)
		if err != nil {
			t.Errorf("Fail to convert %s to ISBN_13: %v", tc.in, err)
			continue
		}

		if tc.out != got {
			t.Errorf("fail to convert %s to ISBN_10. Got: %v. Want: %v", tc.in, got, tc.out)
		}
	}
}

func TestCalcISBN10checkdigit(t *testing.T) {
	testCases := []string{
		"0070350485", "007035958X", "0070393605", "0070418780",
		"0072961945", "007296216X", "0072962216", "0072963050",
		"0323025935", "0324007515", "0324066856", "0805366245",
		"0805368515", "080537146X", "1418843709", "1558606041",
		"155860832X",
	}

	for _, tc := range testCases {
		want := tc[9:]
		got, err := calcISBN10checkdigit(tc[:9])
		if err != nil {
			t.Errorf("Fail to calculate check-digit for %s: %v", tc, err)
			continue
		}

		if want != got {
			t.Errorf("fail to calculate check-digit for %s. Got: %v. Want: %v", tc, got, want)
		}
	}
}

func TestToISBN10(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{in: "9782211238434", out: "2211238432"},
		{in: "9780486285405", out: "0486285405"},
		{in: "9782017064657", out: "2017064653"},
	}

	for _, tc := range testCases {
		got, err := toISBN10(tc.in)
		if err != nil {
			t.Errorf("Fail to convert %s to ISBN_10: %v", tc.in, err)
			continue
		}

		if tc.out != got {
			t.Errorf("fail to convert %s to ISBN_10. Got: %v. Want: %v", tc.in, got, tc.out)
		}
	}
}
