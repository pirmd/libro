package css

import (
	"testing"

	"github.com/gorilla/css/scanner"
)

func TestRuleString(t *testing.T) {
	testCases := []struct {
		in   *Rule
		want string
	}{
		{
			in: &Rule{
				AtKeyword: str2token("@font-face")[0],
				Declarations: []*Declaration{
					{Property: "font-family", Value: str2val("MyHelvetica")},
					{Property: "src", Value: str2val("local(\"Helvetica Neue Bold\"), local(\"HelveticaNeue-Bold\"), url(MgOpenModernaBold.ttf)")},
					{Property: "font-weight", Value: str2val("bold")},
				},
			},
			want: `@font-face {
    font-family: MyHelvetica;
    src: local("Helvetica Neue Bold"), local("HelveticaNeue-Bold"), url(MgOpenModernaBold.ttf);
    font-weight: bold;
}`,
		},
		{
			in: &Rule{
				AtKeyword: str2token("@keyframes")[0],
				Selectors: []Value{str2val("identifier")},
				EmbeddedRuleset: Ruleset{
					&Rule{
						Selectors: []Value{str2val("0%")},
						Declarations: []*Declaration{
							{Property: "top", Value: str2val("0")},
							{Property: "left", Value: str2val("0")},
						},
					},
					&Rule{
						Selectors: []Value{str2val("100%")},
						Declarations: []*Declaration{
							{Property: "top", Value: str2val("100px")},
							{Property: "left", Value: str2val("100%")},
						},
					},
				},
			},
			want: `@keyframes identifier {
    0% {
        top: 0;
        left: 0;
    }
    100% {
        top: 100px;
        left: 100%;
    }
}`,
		},
		{
			in: &Rule{
				AtKeyword: str2token("@supports")[0],
				Selectors: []Value{str2val("(animation-name: test)")},
				EmbeddedRuleset: Ruleset{
					&Rule{
						AtKeyword: str2token("@keyframes")[0],
						EmbeddedRuleset: Ruleset{
							&Rule{
								Selectors: []Value{str2val("0%")},
								Declarations: []*Declaration{
									{Property: "top", Value: str2val("0")},
									{Property: "left", Value: str2val("0")},
								},
							},
							&Rule{
								Selectors: []Value{str2val("100%")},
								Declarations: []*Declaration{
									{Property: "top", Value: str2val("100px")},
									{Property: "left", Value: str2val("100%")},
								},
							},
						},
					},
				},
			},
			want: `@supports (animation-name: test) {
    @keyframes {
        0% {
            top: 0;
            left: 0;
        }
        100% {
            top: 100px;
            left: 100%;
        }
    }
}`,
		},
	}

	for _, tc := range testCases {
		got := tc.in.String()
		if got != tc.want {
			t.Errorf("Formatting ruleset failed.\nWant:\n%s\nGot :\n%s", tc.want, got)
		}
	}
}

func TestDeclarationAppendToValue(t *testing.T) {
	testCases := []struct {
		in        []*scanner.Token
		want      *Declaration
		shallFail bool
	}{
		{
			in:   str2token("blue"),
			want: &Declaration{Value: str2val("blue")},
		},
		{
			in:   str2token("yellow important"),
			want: &Declaration{Value: str2val("yellow important")},
		},
		{
			in:   str2token("green ! important"),
			want: &Declaration{Value: str2val("green ! important")},
		},
		{
			in:   str2token("red !important"),
			want: &Declaration{Value: str2val("red"), IsImportant: true},
		},
		{
			in:        str2token("!important pink"),
			shallFail: true,
		},
	}

NextTC:
	for _, tc := range testCases {
		got := new(Declaration)
		for _, tok := range tc.in {
			if err := got.AppendToValue(tok); err != nil {
				if !tc.shallFail {
					t.Errorf("fail to build Declaration for '%s': %v", tc.in, err)
				}
				continue NextTC
			}
		}

		if tc.shallFail {
			t.Errorf("build Declaration for '%s' does not fail", tc.in)
			continue NextTC
		}

		got.Value = got.Value.Scrub()
		if !areSameDeclaration(tc.want, got) {
			t.Errorf("fail to build Declaration for '%s':\nWant: %s\nGot : %s", tc.in, tc.want, got)
		}
	}
}

func TestValueScrub(t *testing.T) {
	testCases := []struct {
		in        []*scanner.Token
		want      Value
		shallFail bool
	}{
		{
			in:   str2token("yellow !important"),
			want: str2token("yellow !important"),
		},
		{
			in:   str2token(" yellow !important  "),
			want: str2token("yellow !important"),
		},
		{
			in:   str2token("expr/*XSS*/ession(alert('XSS'))"),
			want: str2token("expression(alert('XSS'))"),
		},
	}

	for _, tc := range testCases {
		out := make(Value, 0)
		for _, tok := range tc.in {
			out.Append(tok)
		}
		got := out.Scrub()

		if !areSameValue(tc.want, got) {
			t.Errorf("fail to build Value for '%s':\nWant: %s\nGot : %s", tc.in, tc.want, got)
		}
	}
}
